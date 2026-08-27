package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Device struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	IPAddress  string    `json:"ipAddress"`
	MACAddress string    `json:"macAddress"`
	LEDCount   int       `json:"ledCount"`
	IsOnline   bool      `json:"isOnline"`
	State      WLEDState `json:"state"`
}

type WLEDState struct {
	On         bool          `json:"on"`
	Brightness int           `json:"bri"`
	Transition int           `json:"transition,omitempty"`
	PresetID   int           `json:"ps,omitempty"`
	PlaylistID int           `json:"pl,omitempty"`
	Segments   []WLEDSegment `json:"seg,omitempty"`
}

type WLEDSegment struct {
	ID     int     `json:"id"`
	Start  int     `json:"start"`
	Stop   int     `json:"stop"`
	Length int     `json:"len"`
	Colors [][]int `json:"col,omitempty"`
	FX     int     `json:"fx,omitempty"`
	Speed  int     `json:"sx,omitempty"`
	Ix     int     `json:"ix,omitempty"`
	Pal    int     `json:"pal,omitempty"`
}

type WLEDInfo struct {
	Ver  string `json:"ver"`
	Vid  int    `json:"vid"`
	Leds struct {
		Count int `json:"count"`
		Fps   int `json:"fps"`
		Pwr   int `json:"pwr"`
		Maxpwr int `json:"maxpwr"`
		Maxseg int `json:"maxseg"`
	} `json:"leds"`
	Name string `json:"name"`
	Udp  int    `json:"udp"`
	Mac  string `json:"mac"`
	Ip   string `json:"ip"`
}

type WLEDClient struct {
	httpClient *http.Client
}

func NewWLEDClient() *WLEDClient {
	return &WLEDClient{
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *WLEDClient) FetchDeviceInfo(ip string) (*WLEDInfo, error) {
	url := fmt.Sprintf("http://%s/json/info", ip)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wled returned non-200 status: %d", resp.StatusCode)
	}

	var info WLEDInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *WLEDClient) SetState(ip string, state WLEDState) error {
	url := fmt.Sprintf("http://%s/json/state", ip)
	body, err := json.Marshal(state)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update wled state, status: %d", resp.StatusCode)
	}

	return nil
}

// SendDDPFrame sends raw RGB pixel data over DDP (Distributed Display Protocol) UDP port 4048
func SendDDPFrame(ip string, pixelData []byte) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:4048", ip))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// DDP Header: Flags (0x41 = Push + Version 1), Sequence (0x01), Data Type (0x01 = RGB), ID (0x01)
	header := []byte{
		0x41, 0x01, 0x01, 0x01,
		0x00, 0x00, 0x00, 0x00, // Offset (4 bytes)
		byte((len(pixelData) >> 8) & 0xFF), byte(len(pixelData) & 0xFF), // Data Length (2 bytes)
	}

	packet := append(header, pixelData...)
	_, err = conn.Write(packet)
	return err
}
