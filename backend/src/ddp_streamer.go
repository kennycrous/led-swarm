package main

import (
	"encoding/binary"
	"math"
	"net"
	"sync"
	"time"
)

type DDPStreamTarget struct {
	TargetType string   `json:"targetType"` // "device", "group", "room"
	TargetID   string   `json:"targetID"`
	Effect     string   `json:"effect"`
	Speed      float64  `json:"speed"`
	Intensity  float64  `json:"intensity"`
	IPs        []string `json:"ips"`
	LEDCount   int      `json:"ledCount"`
}

type DDPTargetStatus struct {
	Active     bool    `json:"active"`
	TargetType string  `json:"targetType"`
	TargetID   string  `json:"targetID"`
	Effect     string  `json:"effect"`
	Speed      float64 `json:"speed"`
	Intensity  float64 `json:"intensity"`
	FPS        float64 `json:"fps"`
}

type LEDPixel2D struct {
	IPIndex  int
	LEDIndex int
	X        float64
	Y        float64
}

type activeStream struct {
	target     DDPStreamTarget
	stopChan   chan struct{}
	actualFPS  float64
	frameCount uint64
}

type DDPStreamer struct {
	mu        sync.RWMutex
	streams   map[string]*activeStream
	targetFPS int
	devMgr    *DeviceManager
	groupMgr  *GroupManager
	canvasMgr *CanvasManager
	hub       *Hub
}

func NewDDPStreamer() *DDPStreamer {
	return &DDPStreamer{
		streams:   make(map[string]*activeStream),
		targetFPS: 60,
	}
}

func (ds *DDPStreamer) SetManagers(devMgr *DeviceManager, groupMgr *GroupManager, canvasMgr *CanvasManager, hub *Hub) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.devMgr = devMgr
	ds.groupMgr = groupMgr
	ds.canvasMgr = canvasMgr
	ds.hub = hub
}

func (ds *DDPStreamer) SetHub(hub *Hub) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.hub = hub
}

func (ds *DDPStreamer) SetFPS(fps int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if fps > 0 && fps <= 120 {
		ds.targetFPS = fps
	}
}

func (ds *DDPStreamer) StartStream(targetType string, targetID string, effect string, speed float64, intensity float64, customIPs []string, customLEDCount int) DDPTargetStatus {
	ds.mu.Lock()
	key := targetType + ":" + targetID

	if existing, found := ds.streams[key]; found {
		close(existing.stopChan)
		delete(ds.streams, key)
	}

	if speed <= 0 {
		speed = 1.0
	}
	if intensity <= 0 {
		intensity = 1.0
	}
	if effect == "" {
		effect = "rainbow_wave"
	}

	target := DDPStreamTarget{
		TargetType: targetType,
		TargetID:   targetID,
		Effect:     effect,
		Speed:      speed,
		Intensity:  intensity,
		IPs:        customIPs,
		LEDCount:   customLEDCount,
	}

	stopCh := make(chan struct{})
	stream := &activeStream{
		target:    target,
		stopChan:  stopCh,
		actualFPS: float64(ds.targetFPS),
	}
	ds.streams[key] = stream

	targetFPS := ds.targetFPS
	devMgr := ds.devMgr
	groupMgr := ds.groupMgr
	canvasMgr := ds.canvasMgr
	hub := ds.hub
	ds.mu.Unlock()

	go ds.runStreamLoop(key, target, stopCh, targetFPS, devMgr, groupMgr, canvasMgr, hub)

	return DDPTargetStatus{
		Active:     true,
		TargetType: targetType,
		TargetID:   targetID,
		Effect:     effect,
		Speed:      speed,
		Intensity:  intensity,
		FPS:        float64(targetFPS),
	}
}

func (ds *DDPStreamer) StopStream(targetType string, targetID string) DDPTargetStatus {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	key := targetType + ":" + targetID
	if existing, found := ds.streams[key]; found {
		close(existing.stopChan)
		delete(ds.streams, key)
	}

	status := DDPTargetStatus{
		Active:     false,
		TargetType: targetType,
		TargetID:   targetID,
	}

	if ds.hub != nil {
		ds.hub.BroadcastJSON(map[string]interface{}{
			"type": "ddp_status",
			"data": ds.getAllStatusesLocked(),
		})
	}

	return status
}

func (ds *DDPStreamer) GetStatus() map[string]DDPTargetStatus {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.getAllStatusesLocked()
}

func (ds *DDPStreamer) getAllStatusesLocked() map[string]DDPTargetStatus {
	res := make(map[string]DDPTargetStatus)
	for key, str := range ds.streams {
		res[key] = DDPTargetStatus{
			Active:     true,
			TargetType: str.target.TargetType,
			TargetID:   str.target.TargetID,
			Effect:     str.target.Effect,
			Speed:      str.target.Speed,
			Intensity:  str.target.Intensity,
			FPS:        str.actualFPS,
		}
	}
	return res
}

func (ds *DDPStreamer) runStreamLoop(
	key string,
	target DDPStreamTarget,
	stopCh chan struct{},
	fps int,
	devMgr *DeviceManager,
	groupMgr *GroupManager,
	canvasMgr *CanvasManager,
	hub *Hub,
) {
	tickerInterval := time.Duration(int64(time.Second) / int64(fps))
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	// Resolve target IPs and 2D pixel coordinates
	targetIPs, pixel2DMap, totalLEDs := resolveTargetConfig(target, devMgr, groupMgr, canvasMgr)
	if len(targetIPs) == 0 {
		return
	}

	sockets := make([]net.Conn, 0, len(targetIPs))
	for _, ip := range targetIPs {
		host := ip
		if _, _, err := net.SplitHostPort(ip); err != nil {
			host = net.JoinHostPort(ip, "4048")
		}
		conn, err := net.Dial("udp", host)
		if err == nil {
			sockets = append(sockets, conn)
		}
	}
	defer func() {
		for _, s := range sockets {
			s.Close()
		}
	}()

	startTime := time.Now()
	var seq uint8 = 0
	var framesInWindow uint64 = 0
	lastStatsTime := time.Now()

	for {
		select {
		case <-stopCh:
			return
		case t := <-ticker.C:
			elapsedSec := t.Sub(startTime).Seconds()

			if target.TargetType == "room" && len(pixel2DMap) > 0 {
				// Render 2D Spatial Effects
				perIPBuffers := Generate2DSpatialEffectFrame(target.Effect, pixel2DMap, len(targetIPs), elapsedSec, target.Speed, target.Intensity)
				for ipIdx, buf := range perIPBuffers {
					if ipIdx < len(sockets) && len(buf) > 0 {
						packet := BuildDDPPacket(seq, 0, len(buf), buf)
						_, _ = sockets[ipIdx].Write(packet)
					}
				}
			} else {
				// 1D Linear Buffer Stream
				rgbBuffer := GenerateDDPEffectFrame(target.Effect, totalLEDs, elapsedSec, target.Speed, target.Intensity)
				// Broadcast to all sockets
				packet := BuildDDPPacket(seq, 0, len(rgbBuffer), rgbBuffer)
				for _, s := range sockets {
					_, _ = s.Write(packet)
				}
			}

			seq = (seq + 1) % 16
			framesInWindow++

			if t.Sub(lastStatsTime) >= 500*time.Millisecond {
				dt := t.Sub(lastStatsTime).Seconds()
				currFPS := float64(framesInWindow) / dt

				ds.mu.Lock()
				if str, found := ds.streams[key]; found {
					str.actualFPS = math.Round(currFPS*10) / 10
					str.frameCount += framesInWindow
				}
				allStatuses := ds.getAllStatusesLocked()
				ds.mu.Unlock()

				if hub != nil {
					hub.BroadcastJSON(map[string]interface{}{
						"type": "ddp_status",
						"data": allStatuses,
					})
				}

				framesInWindow = 0
				lastStatsTime = t
			}
		}
	}
}

func resolveTargetConfig(
	target DDPStreamTarget,
	devMgr *DeviceManager,
	groupMgr *GroupManager,
	canvasMgr *CanvasManager,
) ([]string, []LEDPixel2D, int) {
	var ips []string
	var pixels []LEDPixel2D
	totalLEDs := target.LEDCount
	if totalLEDs <= 0 {
		totalLEDs = 60
	}

	if len(target.IPs) > 0 {
		return target.IPs, pixels, totalLEDs
	}

	switch target.TargetType {
	case "device":
		if devMgr != nil {
			dev := devMgr.GetDeviceByID(target.TargetID)
			if dev != nil && dev.IPAddress != "" {
				ips = append(ips, dev.IPAddress)
				if dev.LEDCount > 0 {
					totalLEDs = dev.LEDCount
				}
			}
		}

	case "group":
		if groupMgr != nil {
			group := groupMgr.GetGroupByID(target.TargetID)
			if group != nil {
				for _, devID := range group.DeviceIDs {
					if devMgr != nil {
						dev := devMgr.GetDeviceByID(devID)
						if dev != nil && dev.IPAddress != "" && dev.IsOnline {
							ips = append(ips, dev.IPAddress)
						}
					}
				}
			}
		}

	case "room":
		if canvasMgr != nil {
			placements := canvasMgr.GetPlacementsForRoom(target.TargetID)
			for ipIdx, p := range placements {
				if devMgr != nil {
					dev := devMgr.GetDeviceByID(p.DeviceID)
					if dev != nil && dev.IPAddress != "" && dev.IsOnline {
						ips = append(ips, dev.IPAddress)
						devLEDs := dev.LEDCount
						if devLEDs <= 0 {
							devLEDs = 60
						}
						// Calculate (x, y) coordinates for each LED on 2D grid
						rad := p.Rotation * math.Pi / 180.0
						dirX := math.Cos(rad)
						dirY := math.Sin(rad)
						spacing := (100.0 * p.Scale) / float64(devLEDs)

						for i := 0; i < devLEDs; i++ {
							px := p.PosX + float64(i)*spacing*dirX
							py := p.PosY + float64(i)*spacing*dirY
							pixels = append(pixels, LEDPixel2D{
								IPIndex:  ipIdx,
								LEDIndex: i,
								X:        px,
								Y:        py,
							})
						}
					}
				}
			}
		}
	}

	return ips, pixels, totalLEDs
}

func Generate2DSpatialEffectFrame(
	effect string,
	pixels []LEDPixel2D,
	ipCount int,
	timeSec float64,
	speed float64,
	intensity float64,
) [][]byte {
	buffers := make([][]byte, ipCount)
	t := timeSec * speed

	// Determine max LED index per IP
	maxIndexPerIP := make([]int, ipCount)
	for _, p := range pixels {
		if p.IPIndex < ipCount && p.LEDIndex >= maxIndexPerIP[p.IPIndex] {
			maxIndexPerIP[p.IPIndex] = p.LEDIndex + 1
		}
	}
	for i := 0; i < ipCount; i++ {
		if maxIndexPerIP[i] == 0 {
			maxIndexPerIP[i] = 60
		}
		buffers[i] = make([]byte, maxIndexPerIP[i]*3)
	}

	centerX, centerY := 1000.0, 600.0

	for _, p := range pixels {
		if p.IPIndex >= ipCount {
			continue
		}
		buf := buffers[p.IPIndex]
		idx := p.LEDIndex * 3
		if idx+2 >= len(buf) {
			continue
		}

		var r, g, b uint8

		switch effect {
		case "spatial_sweep":
			sweepPos := math.Mod(t*400.0, 2000.0)
			dist := math.Abs(p.X - sweepPos)
			if dist < 120.0 {
				v := uint8(255 * intensity * (1.0 - dist/120.0))
				r = uint8(float64(v) * 0.2)
				g = v
				b = v
			} else {
				g = uint8(10 * intensity)
			}

		case "spatial_ripple":
			fallthrough
		default:
			dist := math.Sqrt((p.X-centerX)*(p.X-centerX) + (p.Y-centerY)*(p.Y-centerY))
			wave := math.Sin(dist*0.01 - t*3.0)
			if wave > 0 {
				hue := math.Mod(dist*0.001+t*0.1, 1.0)
				r, g, b = hsvToRGB(hue, 1.0, wave*intensity)
			}
		}

		buf[idx] = r
		buf[idx+1] = g
		buf[idx+2] = b
	}

	return buffers
}

func BuildDDPPacket(seq uint8, offset uint32, payloadLength int, payload []byte) []byte {
	packet := make([]byte, 10+payloadLength)
	packet[0] = 0x41 // Push flag + Version 1
	packet[1] = seq
	packet[2] = 0x01 // RGB data type
	packet[3] = 0x01 // Display ID
	binary.BigEndian.PutUint32(packet[4:8], offset)
	binary.BigEndian.PutUint16(packet[8:10], uint16(payloadLength))
	copy(packet[10:], payload)
	return packet
}

func GenerateDDPEffectFrame(effect string, ledCount int, timeSec float64, speed float64, intensity float64) []byte {
	buf := make([]byte, ledCount*3)
	if ledCount <= 0 {
		return buf
	}

	t := timeSec * speed

	switch effect {
	case "digital_rain":
		for i := 0; i < ledCount; i++ {
			pos := float64(i) / float64(ledCount)
			drop := math.Sin((pos*5.0 - t*3.0) * math.Pi)
			if drop > 0.6 {
				v := uint8(255 * intensity * ((drop - 0.6) / 0.4))
				buf[i*3] = v / 4   // R
				buf[i*3+1] = v     // G (Bright green)
				buf[i*3+2] = v / 2 // B
			} else {
				buf[i*3] = 0
				buf[i*3+1] = uint8(20 * intensity)
				buf[i*3+2] = 0
			}
		}

	case "pulse_beads":
		for i := 0; i < ledCount; i++ {
			pos := float64(i) / float64(ledCount)
			bead1 := math.Max(0, 1.0-math.Abs(pos-math.Mod(t*0.8, 1.0))*8.0)
			bead2 := math.Max(0, 1.0-math.Abs(pos-(1.0-math.Mod(t*0.5, 1.0)))*8.0)
			r := uint8(255 * intensity * math.Min(1.0, bead1))
			b := uint8(255 * intensity * math.Min(1.0, bead2))
			g := uint8(100 * intensity * math.Min(1.0, bead1+bead2))

			buf[i*3] = r
			buf[i*3+1] = g
			buf[i*3+2] = b
		}

	case "cyber_fire":
		for i := 0; i < ledCount; i++ {
			pos := float64(i) / float64(ledCount)
			n := math.Sin(pos*8.0+t*4.0)*0.5 + math.Cos(pos*12.0-t*6.0)*0.5
			val := math.Max(0, math.Min(1.0, (n+1.0)/2.0))
			buf[i*3] = uint8(255 * val * intensity)                     // R
			buf[i*3+1] = uint8(80 * (1.0 - val) * intensity)            // G
			buf[i*3+2] = uint8(200 * math.Sin(val*math.Pi) * intensity) // B (Cyber magenta flame)
		}

	case "rainbow_wave":
		fallthrough
	default:
		for i := 0; i < ledCount; i++ {
			hue := math.Mod(float64(i)/float64(ledCount)+t*0.2, 1.0)
			r, g, b := hsvToRGB(hue, 1.0, intensity)
			buf[i*3] = r
			buf[i*3+1] = g
			buf[i*3+2] = b
		}
	}

	return buf
}

func hsvToRGB(h, s, v float64) (uint8, uint8, uint8) {
	i := math.Floor(h * 6.0)
	f := h*6.0 - i
	p := v * (1.0 - s)
	q := v * (1.0 - f*s)
	t := v * (1.0 - (1.0-f)*s)

	var r, g, b float64
	switch int(i) % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}
