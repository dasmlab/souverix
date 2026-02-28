package diagnostics

// ExampleStateProvider shows how components can implement ComponentStateProvider
//
// Example implementation:
//
//	type PcscfState struct {
//		tracker *StateTracker
//		activeCalls map[string]*CallState
//		recordRoutes map[string][]string
//	}
//
//	func (p *PcscfState) GetState(key string) (interface{}, bool) {
//		// Return state from tracker or component's own storage
//		return p.tracker.GetState(key)
//	}
//
//	func (p *PcscfState) SetState(key string, value interface{}) {
//		// Update state in tracker or component's own storage
//		p.tracker.SetState(key, value)
//	}
//
//	func (p *PcscfState) GetAllState() map[string]interface{} {
//		return p.tracker.GetAllState()
//	}
//
//	func (p *PcscfState) GetStateKeys() []string {
//		return p.tracker.GetStateKeys()
//	}
//
//	// In component main.go:
//	// diag := diagnostics.New("Souverix P-CSCF", version, buildTime, gitCommit, logger)
//	// stateProvider := &PcscfState{tracker: diag.GetStateTracker()}
//	// diag.SetStateProvider(stateProvider)
