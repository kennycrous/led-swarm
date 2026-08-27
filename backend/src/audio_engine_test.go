package main

import (
	"testing"
)

func TestAudioEngine_FFTState(t *testing.T) {
	ae := NewAudioEngine()

	if ae.IsActive() {
		t.Errorf("Expected AudioEngine to be inactive initially")
	}

	bins := []float64{0.8, 0.9, 0.7, 0.6, 0.4, 0.3, 0.2, 0.2, 0.1, 0.1, 0.05, 0.05, 0.02, 0.02, 0.01, 0.01}
	ae.UpdateFFT(bins, 0.75, 0.25, 0.03, 0.75)

	if !ae.IsActive() {
		t.Errorf("Expected AudioEngine to be active after updating FFT")
	}

	state := ae.GetAudioState()
	if state.Bass != 0.75 {
		t.Errorf("Expected Bass energy 0.75, got %f", state.Bass)
	}

	if state.Peak != 0.75 {
		t.Errorf("Expected Peak energy 0.75, got %f", state.Peak)
	}
}
