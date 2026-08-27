package main

import (
	"encoding/binary"
	"net"
	"testing"
	"time"
)

func TestDDPHeader_Marshalling(t *testing.T) {
	ledCount := 60
	payloadLength := ledCount * 3
	packet := BuildDDPPacket(1, 0, payloadLength, make([]byte, payloadLength))

	if len(packet) != 10+payloadLength {
		t.Fatalf("Expected packet length %d, got %d", 10+payloadLength, len(packet))
	}

	// Byte 0: Flags (0x41 for Push + v1)
	if packet[0] != 0x41 {
		t.Errorf("Expected flags 0x41, got 0x%02x", packet[0])
	}

	// Byte 1: Sequence
	if packet[1] != 1 {
		t.Errorf("Expected sequence 1, got %d", packet[1])
	}

	// Byte 2: Data type (0x01 = RGB)
	if packet[2] != 0x01 {
		t.Errorf("Expected data type 0x01, got 0x%02x", packet[2])
	}

	// Byte 3: Display ID (0x01)
	if packet[3] != 0x01 {
		t.Errorf("Expected display ID 0x01, got 0x%02x", packet[3])
	}

	// Bytes 4-7: Offset (0)
	offset := binary.BigEndian.Uint32(packet[4:8])
	if offset != 0 {
		t.Errorf("Expected offset 0, got %d", offset)
	}

	// Bytes 8-9: Data length
	dataLen := binary.BigEndian.Uint16(packet[8:10])
	if int(dataLen) != payloadLength {
		t.Errorf("Expected data length %d, got %d", payloadLength, dataLen)
	}
}

func TestDDPEffects_Generators(t *testing.T) {
	ledCount := 30
	effects := []string{"rainbow_wave", "digital_rain", "pulse_beads", "cyber_fire", "audio_bass_pulse", "audio_spectrum_waterfall", "audio_vu_meter", "audio_treble_sparkle"}
	dummyAudio := AudioState{Active: true, Bass: 0.8, Peak: 0.5, Treble: 0.9}

	for _, eff := range effects {
		buf := GenerateDDPEffectFrame(eff, ledCount, 0.5, 1.0, 1.0, dummyAudio)
		if len(buf) != ledCount*3 {
			t.Errorf("Effect %s returned buffer length %d, expected %d", eff, len(buf), ledCount*3)
		}
	}
}

func TestDDPStreamer_Lifecycle(t *testing.T) {
	// Create mock UDP listener on localhost
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to resolve UDP addr: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("Failed to listen UDP: %v", err)
	}
	defer conn.Close()

	destIP := conn.LocalAddr().String()

	streamer := NewDDPStreamer()
	streamer.SetFPS(60)
	streamer.StartStream("device", "dev-1", "rainbow_wave", 1.0, 1.0, []string{destIP}, 30)

	time.Sleep(100 * time.Millisecond)

	statsMap := streamer.GetStatus()
	st, found := statsMap["device:dev-1"]
	if !found || !st.Active {
		t.Errorf("Expected streamer target device:dev-1 to be active")
	}

	// Read UDP packet from mock socket
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	buf := make([]byte, 200)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil || n < 10 {
		t.Errorf("Expected DDP UDP packet received, got error: %v, bytes: %d", err, n)
	}

	streamer.StopStream("device", "dev-1")
	time.Sleep(50 * time.Millisecond)

	statsStop := streamer.GetStatus()
	if _, foundStop := statsStop["device:dev-1"]; foundStop {
		t.Errorf("Expected streamer target device:dev-1 to be inactive after stop")
	}
}
