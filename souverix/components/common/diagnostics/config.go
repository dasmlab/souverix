package diagnostics

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// UnitTestConfig holds the configuration for unit test execution
type UnitTestConfig struct {
	ComponentName    string
	ComponentShort  string
	FlowID          string
	CallerIP        string
	BaseURL         string
	Neighbors       map[string]string // neighbor component -> endpoint URL
	FauxPorts       map[string]int    // component -> faux port
	InternalOps     []string
	FlowDescription string
}

// NewUnitTestConfig creates a unit test configuration from the current context
func (d *Diagnostics) NewUnitTestConfig(flowID, baseURL, callerIP string) (*UnitTestConfig, error) {
	compShortName := extractComponentShortName(d.componentName)
	
	// Get component's participation in this flow
	compMap, exists := d.registry.compMaps[compShortName]
	if !exists {
		return nil, fmt.Errorf("component %s not found in registry", compShortName)
	}
	
	var participation *FlowParticipation
	for _, fp := range compMap.Flows {
		if fp.FlowID == flowID {
			participation = &fp
			break
		}
	}
	
	if participation == nil {
		return nil, fmt.Errorf("component %s does not participate in flow %s", compShortName, flowID)
	}
	
	// Get flow details
	flow, exists := d.registry.GetFlow(flowID)
	if !exists {
		return nil, fmt.Errorf("flow %s not found", flowID)
	}
	
	// Build neighbor endpoints - all point to caller IP (faux server)
	neighbors := make(map[string]string)
	fauxPorts := make(map[string]int)
	basePort := 19000
	
	// Map component names to faux ports
	componentPortMap := map[string]int{
		"pcscf":       basePort,
		"icscf":       basePort + 1,
		"scscf":       basePort + 2,
		"bgcf":        basePort + 3,
		"mgcf":        basePort + 4,
		"hss":         basePort + 5,
		"ue":          basePort + 6,
		"pstn":        basePort + 7,
		"destination": basePort + 8,
	}
	
	// For each neighbor, create endpoint pointing to faux server
	for _, neighbor := range participation.Neighbors {
		neighborLower := strings.ToLower(neighbor)
		// Map component names (e.g., "I-CSCF" -> "icscf")
		if strings.Contains(neighborLower, "cscf") {
			if strings.HasPrefix(neighborLower, "p-") {
				neighborLower = "pcscf"
			} else if strings.HasPrefix(neighborLower, "i-") {
				neighborLower = "icscf"
			} else if strings.HasPrefix(neighborLower, "s-") {
				neighborLower = "scscf"
			}
		}
		
		port, exists := componentPortMap[neighborLower]
		if !exists {
			// Default port for unknown components
			port = basePort + 9
		}
		
		fauxPorts[neighbor] = port
		neighbors[neighbor] = fmt.Sprintf("http://%s:%d/sip", callerIP, port)
	}
	
	config := &UnitTestConfig{
		ComponentName:    d.componentName,
		ComponentShort:   compShortName,
		FlowID:           flowID,
		CallerIP:         callerIP,
		BaseURL:          baseURL,
		Neighbors:        neighbors,
		FauxPorts:        fauxPorts,
		InternalOps:      participation.InternalOps,
		FlowDescription:  flow.Description,
	}
	
	return config, nil
}

// IsDebugEnabled checks if DIAG_DEBUG is enabled
func IsDebugEnabled() bool {
	return os.Getenv("DIAG_DEBUG") == "true" || os.Getenv("DIAG_DEBUG") == "1"
}

// DebugLog logs a message if DIAG_DEBUG is enabled
func DebugLog(logger *logrus.Logger, format string, args ...interface{}) {
	if IsDebugEnabled() {
		logger.Debugf("[DIAG] "+format, args...)
	}
}
