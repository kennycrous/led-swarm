package main

import (
	"sync"
	"time"
)

type AudioState struct {
	Active     bool        `json:"active"`
	Bins       [16]float64 `json:"bins"`
	Bass       float64     `json:"bass"`
	Mid        float64     `json:"mid"`
	Treble     float64     `json:"treble"`
	Peak       float64     `json:"peak"`
	LastUpdate time.Time   `json:"-"`
}

type AudioEngine struct {
	mu    sync.RWMutex
	state AudioState
}

func NewAudioEngine() *AudioEngine {
	return &AudioEngine{
		state: AudioState{
			Active: false,
			Bins:   [16]float64{},
		},
	}
}

func (ae *AudioEngine) UpdateFFT(bins []float64, bass, mid, treble, peak float64) {
	ae.mu.Lock()
	defer ae.mu.Unlock()

	ae.state.Active = true
	ae.state.LastUpdate = time.Now()
	ae.state.Bass = bass
	ae.state.Mid = mid
	ae.state.Treble = treble
	ae.state.Peak = peak

	for i := 0; i < 16 && i < len(bins); i++ {
		ae.state.Bins[i] = bins[i]
	}
}

func (ae *AudioEngine) IsActive() bool {
	ae.mu.RLock()
	defer ae.mu.RUnlock()

	if !ae.state.Active {
		return false
	}
	// Expire audio state if no FFT update received for >2 seconds
	return time.Since(ae.state.LastUpdate) < 2*time.Second
}

func (ae *AudioEngine) GetAudioState() AudioState {
	ae.mu.RLock()
	defer ae.mu.RUnlock()

	stateCopy := ae.state
	if time.Since(ae.state.LastUpdate) > 2*time.Second {
		stateCopy.Active = false
		stateCopy.Bins = [16]float64{}
		stateCopy.Bass = 0
		stateCopy.Mid = 0
		stateCopy.Treble = 0
		stateCopy.Peak = 0
	}
	return stateCopy
}
