package diagnostics

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// StateVerifier verifies component state changes during call flow execution
type StateVerifier struct {
	componentName string
	stateStore    map[string]interface{} // Component's state store
	verifications []VerificationResult
}

// NewStateVerifier creates a new state verifier
func NewStateVerifier(componentName string) *StateVerifier {
	return &StateVerifier{
		componentName: componentName,
		stateStore:    make(map[string]interface{}),
		verifications: []VerificationResult{},
	}
}

// VerificationResult represents the result of a state verification
type VerificationResult struct {
	Step        int       // Call flow step number
	Operation   string    // Operation being verified
	Key         string    // State key being checked
	Expected    interface{} // Expected value
	Actual      interface{} // Actual value
	Passed      bool      // Whether verification passed
	Message     string    // Verification message
	Timestamp   time.Time // When verification occurred
}

// StateChange represents a change to component state
type StateChange struct {
	Key       string      // State key
	Value     interface{} // New value
	Operation string     // Operation type (set, update, delete, etc.)
	Timestamp time.Time  // When change occurred
}

// RecordStateChange records a state change in the component
func (v *StateVerifier) RecordStateChange(change StateChange) {
	v.stateStore[change.Key] = change.Value
}

// VerifyState verifies that a state key has the expected value
func (v *StateVerifier) VerifyState(step int, operation, key string, expected interface{}) VerificationResult {
	actual, exists := v.stateStore[key]
	
	result := VerificationResult{
		Step:      step,
		Operation: operation,
		Key:       key,
		Expected:  expected,
		Actual:    actual,
		Timestamp: time.Now(),
	}

	if !exists {
		result.Passed = false
		result.Message = fmt.Sprintf("State key '%s' does not exist", key)
	} else if !reflect.DeepEqual(expected, actual) {
		result.Passed = false
		result.Message = fmt.Sprintf("State key '%s' mismatch: expected %v, got %v", key, expected, actual)
	} else {
		result.Passed = true
		result.Message = fmt.Sprintf("State key '%s' verified: %v", key, actual)
	}

	v.verifications = append(v.verifications, result)
	return result
}

// VerifyStateExists verifies that a state key exists
func (v *StateVerifier) VerifyStateExists(step int, operation, key string) VerificationResult {
	_, exists := v.stateStore[key]
	
	result := VerificationResult{
		Step:      step,
		Operation: operation,
		Key:       key,
		Timestamp: time.Now(),
	}

	if exists {
		result.Passed = true
		result.Message = fmt.Sprintf("State key '%s' exists", key)
	} else {
		result.Passed = false
		result.Message = fmt.Sprintf("State key '%s' does not exist", key)
	}

	v.verifications = append(v.verifications, result)
	return result
}

// VerifyStateNotExists verifies that a state key does not exist
func (v *StateVerifier) VerifyStateNotExists(step int, operation, key string) VerificationResult {
	_, exists := v.stateStore[key]
	
	result := VerificationResult{
		Step:      step,
		Operation: operation,
		Key:       key,
		Timestamp: time.Now(),
	}

	if !exists {
		result.Passed = true
		result.Message = fmt.Sprintf("State key '%s' does not exist (as expected)", key)
	} else {
		result.Passed = false
		result.Message = fmt.Sprintf("State key '%s' exists (should not)", key)
	}

	v.verifications = append(v.verifications, result)
	return result
}

// GetVerifications returns all verification results
func (v *StateVerifier) GetVerifications() []VerificationResult {
	return v.verifications
}

// GetPassedVerifications returns only passed verifications
func (v *StateVerifier) GetPassedVerifications() []VerificationResult {
	var passed []VerificationResult
	for _, v := range v.verifications {
		if v.Passed {
			passed = append(passed, v)
		}
	}
	return passed
}

// GetFailedVerifications returns only failed verifications
func (v *StateVerifier) GetFailedVerifications() []VerificationResult {
	var failed []VerificationResult
	for _, v := range v.verifications {
		if !v.Passed {
			failed = append(failed, v)
		}
	}
	return failed
}

// AllPassed returns true if all verifications passed
func (v *StateVerifier) AllPassed() bool {
	for _, v := range v.verifications {
		if !v.Passed {
			return false
		}
	}
	return true
}

// GetSummary returns a summary of verification results
func (v *StateVerifier) GetSummary() map[string]interface{} {
	total := len(v.verifications)
	passed := len(v.GetPassedVerifications())
	failed := len(v.GetFailedVerifications())

	return map[string]interface{}{
		"component":    v.componentName,
		"total":         total,
		"passed":        passed,
		"failed":        failed,
		"all_passed":    v.AllPassed(),
		"verifications": v.verifications,
	}
}

// ComponentStateProvider is an interface that components can implement
// to provide their state for verification
type ComponentStateProvider interface {
	GetState(key string) (interface{}, bool)
	SetState(key string, value interface{})
	GetAllState() map[string]interface{}
	GetStateKeys() []string
}

// VerifyComponentState verifies component state using the provider interface
func (v *StateVerifier) VerifyComponentState(step int, operation string, provider ComponentStateProvider, key string, expected interface{}) VerificationResult {
	actual, exists := provider.GetState(key)
	
	result := VerificationResult{
		Step:      step,
		Operation: operation,
		Key:       key,
		Expected:  expected,
		Actual:    actual,
		Timestamp: time.Now(),
	}

	if !exists {
		result.Passed = false
		result.Message = fmt.Sprintf("State key '%s' does not exist in component", key)
	} else if !reflect.DeepEqual(expected, actual) {
		result.Passed = false
		result.Message = fmt.Sprintf("State key '%s' mismatch: expected %v, got %v", key, expected, actual)
	} else {
		result.Passed = true
		result.Message = fmt.Sprintf("State key '%s' verified: %v", key, actual)
	}

	v.verifications = append(v.verifications, result)
	return result
}

// ExportVerifications exports verification results as JSON
func (v *StateVerifier) ExportVerifications() ([]byte, error) {
	summary := v.GetSummary()
	return json.MarshalIndent(summary, "", "  ")
}
