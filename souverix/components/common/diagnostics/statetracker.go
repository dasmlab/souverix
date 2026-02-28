package diagnostics

import (
	"sync"
	"time"
)

// StateTracker tracks component state changes during call flows
type StateTracker struct {
	mu         sync.RWMutex
	state      map[string]interface{}
	history    []StateChange
	callStates map[string]*CallState // Call-ID -> CallState
}

// CallState represents the state of a SIP call/dialog
type CallState struct {
	CallID      string
	From        string
	To          string
	Method      string
	Status      string // "initiated", "ringing", "answered", "terminated"
	StartTime   time.Time
	EndTime     *time.Time
	Headers     map[string]string
	Body        string
	RouteSet    []string
	RecordRoute []string
	State       map[string]interface{} // Component-specific state
}

// NewStateTracker creates a new state tracker
func NewStateTracker() *StateTracker {
	return &StateTracker{
		state:      make(map[string]interface{}),
		history:    []StateChange{},
		callStates: make(map[string]*CallState),
	}
}

// SetState sets a state value
func (st *StateTracker) SetState(key string, value interface{}) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.state[key] = value
	st.history = append(st.history, StateChange{
		Key:       key,
		Value:     value,
		Operation: "set",
		Timestamp: time.Now(),
	})
}

// GetState gets a state value
func (st *StateTracker) GetState(key string) (interface{}, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	value, exists := st.state[key]
	return value, exists
}

// DeleteState deletes a state value
func (st *StateTracker) DeleteState(key string) {
	st.mu.Lock()
	defer st.mu.Unlock()

	delete(st.state, key)
	st.history = append(st.history, StateChange{
		Key:       key,
		Value:     nil,
		Operation: "delete",
		Timestamp: time.Now(),
	})
}

// GetAllState returns all state
func (st *StateTracker) GetAllState() map[string]interface{} {
	st.mu.RLock()
	defer st.mu.RUnlock()

	result := make(map[string]interface{})
	for k, v := range st.state {
		result[k] = v
	}
	return result
}

// GetStateKeys returns all state keys
func (st *StateTracker) GetStateKeys() []string {
	st.mu.RLock()
	defer st.mu.RUnlock()

	keys := make([]string, 0, len(st.state))
	for k := range st.state {
		keys = append(keys, k)
	}
	return keys
}

// GetHistory returns state change history
func (st *StateTracker) GetHistory() []StateChange {
	st.mu.RLock()
	defer st.mu.RUnlock()

	result := make([]StateChange, len(st.history))
	copy(result, st.history)
	return result
}

// CreateCallState creates a new call state
func (st *StateTracker) CreateCallState(callID, from, to, method string) *CallState {
	st.mu.Lock()
	defer st.mu.Unlock()

	callState := &CallState{
		CallID:      callID,
		From:        from,
		To:          to,
		Method:      method,
		Status:      "initiated",
		StartTime:   time.Now(),
		Headers:     make(map[string]string),
		RouteSet:    []string{},
		RecordRoute: []string{},
		State:       make(map[string]interface{}),
	}

	st.callStates[callID] = callState
	return callState
}

// GetCallState gets a call state by Call-ID
func (st *StateTracker) GetCallState(callID string) (*CallState, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	callState, exists := st.callStates[callID]
	return callState, exists
}

// UpdateCallState updates a call state
func (st *StateTracker) UpdateCallState(callID string, updater func(*CallState)) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if callState, exists := st.callStates[callID]; exists {
		updater(callState)
	}
}

// SetCallStatus sets the status of a call
func (st *StateTracker) SetCallStatus(callID, status string) {
	st.UpdateCallState(callID, func(cs *CallState) {
		cs.Status = status
		if status == "terminated" {
			now := time.Now()
			cs.EndTime = &now
		}
	})
}

// AddRecordRoute adds a Record-Route header to a call
func (st *StateTracker) AddRecordRoute(callID, route string) {
	st.UpdateCallState(callID, func(cs *CallState) {
		cs.RecordRoute = append(cs.RecordRoute, route)
	})
}

// AddRoute adds a Route header to a call
func (st *StateTracker) AddRoute(callID, route string) {
	st.UpdateCallState(callID, func(cs *CallState) {
		cs.RouteSet = append(cs.RouteSet, route)
	})
}

// SetCallHeader sets a header for a call
func (st *StateTracker) SetCallHeader(callID, header, value string) {
	st.UpdateCallState(callID, func(cs *CallState) {
		cs.Headers[header] = value
	})
}

// SetCallStateValue sets a component-specific state value for a call
func (st *StateTracker) SetCallStateValue(callID, key string, value interface{}) {
	st.UpdateCallState(callID, func(cs *CallState) {
		cs.State[key] = value
	})
}

// GetAllCallStates returns all call states
func (st *StateTracker) GetAllCallStates() map[string]*CallState {
	st.mu.RLock()
	defer st.mu.RUnlock()

	result := make(map[string]*CallState)
	for k, v := range st.callStates {
		result[k] = v
	}
	return result
}

// Clear clears all state
func (st *StateTracker) Clear() {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.state = make(map[string]interface{})
	st.history = []StateChange{}
	st.callStates = make(map[string]*CallState)
}
