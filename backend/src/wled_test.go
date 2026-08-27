package main

import (
	"encoding/json"
	"testing"
)

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
