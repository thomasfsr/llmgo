package main

import (
	"sync"
)

type OverallState struct {
	UserID       int
	UserInput    string
	Messages     []string
	ExerciseList ListOfExercises
}

type StateManager struct {
	mu     sync.RWMutex
	states map[string]*OverallState
}

func NewStateManager() *StateManager {
	return &StateManager{states: make(map[string]*OverallState)}
}

func (m *StateManager) Get(threadID string) *OverallState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.states[threadID]
}

func (m *StateManager) Set(threadID string, s *OverallState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[threadID] = s
}
