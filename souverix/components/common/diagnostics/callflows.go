package diagnostics

// CallFlow represents a standard call flow from ETSI/3GPP specifications
type CallFlow struct {
	ID          string   // Unique identifier (e.g., "IMS_REGISTER_AKA", "SIP_INVITE_IMS_TO_IMS", "SIP_INVITE_IMS_TO_PSTN")
	Name        string   // Human-readable name
	Description string   // Description of the call flow
	Spec        string   // ETSI/3GPP specification reference
	Steps       []Step   // Sequence of steps in the call flow
}

// Step represents a single step in a call flow
type Step struct {
	Sequence    int      // Step sequence number
	From        string   // Source component/node
	To          string   // Destination component/node
	Message     string   // Message type (e.g., "INVITE", "REGISTER", "200 OK")
	Interface   string   // Interface name (e.g., "Gm", "Mw", "Mi", "Mj")
	Direction   string   // "request" or "response"
	Description string   // What happens in this step
}

// ComponentCallFlowMap maps components to the call flows they participate in
// and which steps they are involved in
type ComponentCallFlowMap struct {
	Component string   // Component name (e.g., "pcscf", "icscf", "scscf")
	Flows     []FlowParticipation
}

// FlowParticipation describes how a component participates in a call flow
type FlowParticipation struct {
	FlowID      string   // Call flow ID
	Steps       []int    // Step sequence numbers this component participates in
	Role        string   // Component's role in this flow (e.g., "proxy", "handler", "interrogator")
	Neighbors   []string // 1-hop neighbors in this flow
	InternalOps []string // Internal operations this component performs
}

// CallFlowRegistry maintains the master list of all call flows
type CallFlowRegistry struct {
	flows    map[string]*CallFlow
	compMaps map[string]*ComponentCallFlowMap
}

// NewCallFlowRegistry creates a new call flow registry
func NewCallFlowRegistry() *CallFlowRegistry {
	registry := &CallFlowRegistry{
		flows:    make(map[string]*CallFlow),
		compMaps: make(map[string]*ComponentCallFlowMap),
	}
	registry.initStandardFlows()
	return registry
}

// initStandardFlows initializes standard ETSI/3GPP call flows
func (r *CallFlowRegistry) initStandardFlows() {
	// SIP REGISTER with IMS AKA
	r.flows["IMS_REGISTER_AKA"] = &CallFlow{
		ID:          "IMS_REGISTER_AKA",
		Name:        "IMS Registration with AKA Authentication",
		Description: "UE registers with IMS network using AKA authentication",
		Spec:        "3GPP TS 24.229, 3GPP TS 33.203",
		Steps: []Step{
			{Sequence: 1, From: "UE", To: "P-CSCF", Message: "REGISTER", Interface: "Gm", Direction: "request", Description: "UE sends REGISTER request"},
			{Sequence: 2, From: "P-CSCF", To: "I-CSCF", Message: "REGISTER", Interface: "Mw", Direction: "request", Description: "P-CSCF forwards REGISTER"},
			{Sequence: 3, From: "I-CSCF", To: "HSS", Message: "UAR", Interface: "Cx", Direction: "request", Description: "I-CSCF queries HSS for S-CSCF assignment"},
			{Sequence: 4, From: "HSS", To: "I-CSCF", Message: "UAA", Interface: "Cx", Direction: "response", Description: "HSS returns S-CSCF assignment"},
			{Sequence: 5, From: "I-CSCF", To: "S-CSCF", Message: "REGISTER", Interface: "Mw", Direction: "request", Description: "I-CSCF forwards to S-CSCF"},
			{Sequence: 6, From: "S-CSCF", To: "UE", Message: "401 Unauthorized", Interface: "Mw", Direction: "response", Description: "S-CSCF challenges with AKA"},
			{Sequence: 7, From: "UE", To: "P-CSCF", Message: "REGISTER + Authorization", Interface: "Gm", Direction: "request", Description: "UE responds with AKA response"},
			{Sequence: 8, From: "P-CSCF", To: "I-CSCF", Message: "REGISTER + Authorization", Interface: "Mw", Direction: "request", Description: "P-CSCF forwards"},
			{Sequence: 9, From: "I-CSCF", To: "S-CSCF", Message: "REGISTER + Authorization", Interface: "Mw", Direction: "request", Description: "I-CSCF forwards"},
			{Sequence: 10, From: "S-CSCF", To: "HSS", Message: "SAR", Interface: "Cx", Direction: "request", Description: "S-CSCF requests service profile"},
			{Sequence: 11, From: "HSS", To: "S-CSCF", Message: "SAA", Interface: "Cx", Direction: "response", Description: "HSS returns service profile"},
			{Sequence: 12, From: "S-CSCF", To: "UE", Message: "200 OK", Interface: "Mw", Direction: "response", Description: "Registration successful"},
		},
	}

	// SIP INVITE IMS-to-IMS
	r.flows["SIP_INVITE_IMS_TO_IMS"] = &CallFlow{
		ID:          "SIP_INVITE_IMS_TO_IMS",
		Name:        "SIP INVITE - IMS to IMS",
		Description: "Basic SIP INVITE call between two IMS users",
		Spec:        "3GPP TS 24.229",
		Steps: []Step{
			{Sequence: 1, From: "UE", To: "P-CSCF", Message: "INVITE", Interface: "Gm", Direction: "request", Description: "UE sends INVITE"},
			{Sequence: 2, From: "P-CSCF", To: "I-CSCF", Message: "INVITE", Interface: "Mw", Direction: "request", Description: "P-CSCF forwards"},
			{Sequence: 3, From: "I-CSCF", To: "S-CSCF", Message: "INVITE", Interface: "Mw", Direction: "request", Description: "I-CSCF forwards"},
			{Sequence: 4, From: "S-CSCF", To: "Destination", Message: "INVITE", Interface: "ISC", Direction: "request", Description: "S-CSCF routes to destination"},
			{Sequence: 5, From: "Destination", To: "S-CSCF", Message: "180 Ringing", Interface: "ISC", Direction: "response", Description: "Destination alerts"},
			{Sequence: 6, From: "S-CSCF", To: "UE", Message: "180 Ringing", Interface: "Mw", Direction: "response", Description: "Ringing indication"},
			{Sequence: 7, From: "Destination", To: "S-CSCF", Message: "200 OK", Interface: "ISC", Direction: "response", Description: "Call answered"},
			{Sequence: 8, From: "S-CSCF", To: "UE", Message: "200 OK", Interface: "Mw", Direction: "response", Description: "Call established"},
		},
	}

	// SIP INVITE IMS-to-PSTN (VoLTE)
	r.flows["SIP_INVITE_IMS_TO_PSTN"] = &CallFlow{
		ID:          "SIP_INVITE_IMS_TO_PSTN",
		Name:        "SIP INVITE - IMS to PSTN (VoLTE)",
		Description: "SIP INVITE call from IMS to PSTN network",
		Spec:        "3GPP TS 29.163",
		Steps: []Step{
			{Sequence: 1, From: "UE", To: "P-CSCF", Message: "INVITE", Interface: "Gm", Direction: "request", Description: "UE sends INVITE with tel: URI"},
			{Sequence: 2, From: "P-CSCF", To: "I-CSCF", Message: "INVITE", Interface: "Mw", Direction: "request", Description: "P-CSCF forwards"},
			{Sequence: 3, From: "I-CSCF", To: "S-CSCF", Message: "INVITE", Interface: "Mw", Direction: "request", Description: "I-CSCF forwards"},
			{Sequence: 4, From: "S-CSCF", To: "BGCF", Message: "INVITE", Interface: "Mi", Direction: "request", Description: "S-CSCF detects PSTN breakout needed"},
			{Sequence: 5, From: "BGCF", To: "MGCF", Message: "INVITE", Interface: "Mj", Direction: "request", Description: "BGCF selects MGCF"},
			{Sequence: 6, From: "MGCF", To: "PSTN", Message: "ISUP IAM", Interface: "Nc", Direction: "request", Description: "MGCF converts to ISUP"},
			{Sequence: 7, From: "PSTN", To: "MGCF", Message: "ISUP ACM", Interface: "Nc", Direction: "response", Description: "PSTN alerts"},
			{Sequence: 8, From: "MGCF", To: "BGCF", Message: "180 Ringing", Interface: "Mj", Direction: "response", Description: "MGCF converts to SIP"},
			{Sequence: 9, From: "BGCF", To: "S-CSCF", Message: "180 Ringing", Interface: "Mi", Direction: "response", Description: "BGCF forwards"},
			{Sequence: 10, From: "S-CSCF", To: "UE", Message: "180 Ringing", Interface: "Mw", Direction: "response", Description: "Ringing indication"},
			{Sequence: 11, From: "PSTN", To: "MGCF", Message: "ISUP ANM", Interface: "Nc", Direction: "response", Description: "PSTN answers"},
			{Sequence: 12, From: "MGCF", To: "BGCF", Message: "200 OK", Interface: "Mj", Direction: "response", Description: "MGCF converts to SIP"},
			{Sequence: 13, From: "BGCF", To: "S-CSCF", Message: "200 OK", Interface: "Mi", Direction: "response", Description: "BGCF forwards"},
			{Sequence: 14, From: "S-CSCF", To: "UE", Message: "200 OK", Interface: "Mw", Direction: "response", Description: "Call established"},
		},
	}

	// Initialize component maps
	r.initComponentMaps()
}

// initComponentMaps initializes component participation maps
func (r *CallFlowRegistry) initComponentMaps() {
	// P-CSCF participation
	r.compMaps["pcscf"] = &ComponentCallFlowMap{
		Component: "pcscf",
		Flows: []FlowParticipation{
			{
				FlowID:    "IMS_REGISTER_AKA",
				Steps:     []int{1, 2, 7, 8},
				Role:      "proxy",
				Neighbors: []string{"UE", "I-CSCF"},
				InternalOps: []string{"validate_headers", "insert_record_route", "security_validation"},
			},
			{
				FlowID:    "SIP_INVITE_IMS_TO_IMS",
				Steps:     []int{1, 2},
				Role:      "proxy",
				Neighbors: []string{"UE", "I-CSCF"},
				InternalOps: []string{"validate_headers", "insert_record_route", "policy_control"},
			},
			{
				FlowID:    "SIP_INVITE_IMS_TO_PSTN",
				Steps:     []int{1, 2},
				Role:      "proxy",
				Neighbors: []string{"UE", "I-CSCF"},
				InternalOps: []string{"validate_headers", "insert_record_route", "policy_control"},
			},
		},
	}

	// I-CSCF participation
	r.compMaps["icscf"] = &ComponentCallFlowMap{
		Component: "icscf",
		Flows: []FlowParticipation{
			{
				FlowID:    "IMS_REGISTER_AKA",
				Steps:     []int{2, 3, 4, 5, 8, 9},
				Role:      "interrogator",
				Neighbors: []string{"P-CSCF", "HSS", "S-CSCF"},
				InternalOps: []string{"query_hss_for_scscf", "route_to_scscf"},
			},
			{
				FlowID:    "SIP_INVITE_IMS_TO_IMS",
				Steps:     []int{2, 3},
				Role:      "interrogator",
				Neighbors: []string{"P-CSCF", "S-CSCF"},
				InternalOps: []string{"query_hss_for_scscf", "route_to_scscf"},
			},
			{
				FlowID:    "SIP_INVITE_IMS_TO_PSTN",
				Steps:     []int{2, 3},
				Role:      "interrogator",
				Neighbors: []string{"P-CSCF", "S-CSCF"},
				InternalOps: []string{"query_hss_for_scscf", "route_to_scscf"},
			},
		},
	}

	// S-CSCF participation
	r.compMaps["scscf"] = &ComponentCallFlowMap{
		Component: "scscf",
		Flows: []FlowParticipation{
			{
				FlowID:    "IMS_REGISTER_AKA",
				Steps:     []int{5, 6, 9, 10, 11, 12},
				Role:      "handler",
				Neighbors: []string{"I-CSCF", "HSS"},
				InternalOps: []string{"aka_challenge", "verify_aka_response", "load_service_profile", "apply_ifc"},
			},
			{
				FlowID:    "SIP_INVITE_IMS_TO_IMS",
				Steps:     []int{3, 4, 5, 6, 7, 8},
				Role:      "handler",
				Neighbors: []string{"I-CSCF", "Destination"},
				InternalOps: []string{"load_service_profile", "apply_ifc", "route_to_destination"},
			},
			{
				FlowID:    "SIP_INVITE_IMS_TO_PSTN",
				Steps:     []int{3, 4},
				Role:      "handler",
				Neighbors: []string{"I-CSCF", "BGCF"},
				InternalOps: []string{"detect_pstn_breakout", "route_to_bgcf"},
			},
		},
	}

	// BGCF participation
	r.compMaps["bgcf"] = &ComponentCallFlowMap{
		Component: "bgcf",
		Flows: []FlowParticipation{
			{
				FlowID:    "SIP_INVITE_IMS_TO_PSTN",
				Steps:     []int{4, 5, 8, 9, 12, 13},
				Role:      "breakout_selector",
				Neighbors: []string{"S-CSCF", "MGCF"},
				InternalOps: []string{"select_breakout_network", "select_mgcf", "generate_cdr"},
			},
		},
	}

	// MGCF participation
	r.compMaps["mgcf"] = &ComponentCallFlowMap{
		Component: "mgcf",
		Flows: []FlowParticipation{
			{
				FlowID:    "SIP_INVITE_IMS_TO_PSTN",
				Steps:     []int{5, 6, 7, 8, 11, 12},
				Role:      "interworking",
				Neighbors: []string{"BGCF", "PSTN"},
				InternalOps: []string{"sip_to_isup_conversion", "control_mgw", "isup_to_sip_conversion"},
			},
		},
	}

	// HSS participation
	r.compMaps["hss"] = &ComponentCallFlowMap{
		Component: "hss",
		Flows: []FlowParticipation{
			{
				FlowID:    "IMS_REGISTER_AKA",
				Steps:     []int{3, 4, 10, 11},
				Role:      "database",
				Neighbors: []string{"I-CSCF", "S-CSCF"},
				InternalOps: []string{"assign_scscf", "return_service_profile", "store_registration_state"},
			},
		},
	}
}

// GetFlow returns a call flow by ID
func (r *CallFlowRegistry) GetFlow(flowID string) (*CallFlow, bool) {
	flow, exists := r.flows[flowID]
	return flow, exists
}

// GetComponentFlows returns all call flows a component participates in
func (r *CallFlowRegistry) GetComponentFlows(component string) []FlowParticipation {
	if compMap, exists := r.compMaps[component]; exists {
		return compMap.Flows
	}
	return []FlowParticipation{}
}

// GetComponentSteps returns the steps a component participates in for a given flow
func (r *CallFlowRegistry) GetComponentSteps(component, flowID string) []Step {
	compMap, exists := r.compMaps[component]
	if !exists {
		return []Step{}
	}

	flow, exists := r.flows[flowID]
	if !exists {
		return []Step{}
	}

	for _, participation := range compMap.Flows {
		if participation.FlowID == flowID {
			var steps []Step
			for _, seqNum := range participation.Steps {
				for _, step := range flow.Steps {
					if step.Sequence == seqNum {
						steps = append(steps, step)
						break
					}
				}
			}
			return steps
		}
	}

	return []Step{}
}
