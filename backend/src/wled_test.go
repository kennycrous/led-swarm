package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createMockWLEDServer(t *testing.T) (*httptest.Server, *WLEDState) {
	t.Helper()

	state := &WLEDState{
		On:         true,
		Brightness: 180,
		Segments: []WLEDSegment{
			{
				ID:     0,
				Start:  0,
				Stop:   60,
				Length: 60,
				FX:     0,
				Pal:    0,
				Colors: [][]int{{6, 182, 212}},
			},
		},
	}

	mux := http.NewServeMux()

	// 1. Mock /json/info endpoint
	mux.HandleFunc("/json/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ver":  "0.14.0",
			"name": "Mock WLED Strip",
			"mac":  "AA:BB:CC:DD:EE:FF",
			"leds": map[string]interface{}{
				"count": 60,
			},
		})
	})

	// 2. Mock /json/state endpoint
	mux.HandleFunc("/json/state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			var incoming WLEDState
			json.NewDecoder(r.Body).Decode(&incoming)
			if incoming.On != state.On {
				state.On = incoming.On
			}
			if incoming.Brightness > 0 {
				state.Brightness = incoming.Brightness
			}
			if len(incoming.Segments) > 0 {
				state.Segments = incoming.Segments
			}
		}
		json.NewEncoder(w).Encode(state)
	})

	// 3. Mock /json/eff and /json/pal endpoints
	mux.HandleFunc("/json/eff", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{"Solid", "Breathe", "Rainbow", "Fire 2012"})
	})
	mux.HandleFunc("/json/pal", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{"Default", "Random Cycle", "Colorwaves"})
	})

	ts := httptest.NewServer(mux)
	t.Cleanup(func() {
		ts.Close()
	})

	return ts, state
}

func TestWLEDClient_DeviceInfoAndState(t *testing.T) {
	ts, state := createMockWLEDServer(t)
	client := NewWLEDClient()

	// Extract IP from test server URL (http://127.0.0.1:port)
	host := strings.TrimPrefix(ts.URL, "http://")

	// 1. FetchDeviceInfo
	info, err := client.FetchDeviceInfo(host)
	if err != nil {
		t.Fatalf("FetchDeviceInfo failed: %v", err)
	}
	if info.Name != "Mock WLED Strip" {
		t.Errorf("Expected Name 'Mock WLED Strip', got '%s'", info.Name)
	}
	if info.Leds.Count != 60 {
		t.Errorf("Expected LEDCount 60, got %d", info.Leds.Count)
	}
	if info.Mac != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("Expected MAC 'AA:BB:CC:DD:EE:FF', got '%s'", info.Mac)
	}

	// 2. FetchLiveState
	rawState, err := client.FetchLiveState(host)
	if err != nil {
		t.Fatalf("FetchLiveState failed: %v", err)
	}
	var currState WLEDState
	json.Unmarshal(rawState, &currState)
	if !currState.On {
		t.Errorf("Expected On=true")
	}
	if currState.Brightness != 180 {
		t.Errorf("Expected Brightness 180, got %d", currState.Brightness)
	}

	// 3. SetState (Turn off)
	newState := WLEDState{On: false, Brightness: 50}
	if err := client.SetState(host, newState); err != nil {
		t.Fatalf("SetState failed: %v", err)
	}

	if state.On != false {
		t.Errorf("Expected mock server state.On=false after SetState")
	}
}

func TestWLEDState_JSONMarshalling(t *testing.T) {
	state := WLEDState{
		On:         true,
		Brightness: 200,
		Segments: []WLEDSegment{
			{
				ID:    0,
				FX:    2,
				Speed: 128,
				Pal:   5,
				Colors: [][]int{
					{255, 0, 128},
				},
			},
		},
	}

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal WLEDState: %v", err)
	}

	var restored WLEDState
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Failed to unmarshal WLEDState: %v", err)
	}

	if !restored.On {
		t.Errorf("Expected On=true")
	}
	if restored.Brightness != 200 {
		t.Errorf("Expected Brightness=200, got %d", restored.Brightness)
	}
	if len(restored.Segments) != 1 {
		t.Fatalf("Expected 1 segment")
	}
	if restored.Segments[0].FX != 2 {
		t.Errorf("Expected FX=2, got %d", restored.Segments[0].FX)
	}
}
