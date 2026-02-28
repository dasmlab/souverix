package diagnostics

// ExampleUsage demonstrates how components should use the diagnostics framework
// This file serves as documentation and can be removed in production

/*
Example: Component Integration

1. Initialize diagnostics in component main.go:

	import "github.com/dasmlab/souverix/common/diagnostics"

	func main() {
		logger := logrus.New()
		diag := diagnostics.New("Souverix P-CSCF", version, buildTime, gitCommit, logger)
		
		// Register routes on diagnostics server (r3)
		diag.RegisterRoutes(r3)
	}

2. Implement ComponentStateProvider interface for state verification:

	type PcscfState struct {
		activeCalls map[string]*CallState
		recordRoutes map[string]string
	}

	func (p *PcscfState) GetState(key string) (interface{}, bool) {
		// Return component state
	}

	func (p *PcscfState) SetState(key string, value interface{}) {
		// Update component state
	}

	func (p *PcscfState) GetAllState() map[string]interface{} {
		// Return all state
	}

	func (p *PcscfState) GetStateKeys() []string {
		// Return all state keys
	}

3. Use unit test endpoint:

	GET /diag/unit_test?flow_id=SIP_INVITE_IMS_TO_IMS&base_url=http://localhost:8081

	Response:
	{
		"component": "pcscf",
		"flow_id": "SIP_INVITE_IMS_TO_IMS",
		"all_passed": true,
		"steps": [
			{
				"step": 1,
				"message": "INVITE",
				"interface": "Gm",
				"direction": "request",
				"request": {...},
				"response": {...},
				"passed": true,
				"state_verification": {...}
			}
		],
		"verification": {
			"total": 2,
			"passed": 2,
			"failed": 0,
			"all_passed": true
		}
	}

4. Runtime unit testing via Go subroutines:

	// Component can trigger unit tests during runtime
	go func() {
		// Call own /diag/unit_test endpoint
		resp, err := http.Get("http://localhost:9081/diag/unit_test?flow_id=SIP_INVITE_IMS_TO_IMS")
		// Process results
	}()
*/
