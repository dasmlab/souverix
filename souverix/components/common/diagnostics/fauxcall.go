package diagnostics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FauxCallGenerator generates fake requests and responses for call flow simulation
type FauxCallGenerator struct {
	registry *CallFlowRegistry
	client   *http.Client
}

// NewFauxCallGenerator creates a new faux call generator
func NewFauxCallGenerator(registry *CallFlowRegistry) *FauxCallGenerator {
	return &FauxCallGenerator{
		registry: registry,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// FauxRequest represents a fake request to be sent
type FauxRequest struct {
	Method      string            // HTTP method (e.g., "POST", "GET")
	URL         string            // Target URL
	Headers     map[string]string // Request headers
	Body        string            // Request body (SIP message, JSON, etc.)
	ExpectedCode int             // Expected HTTP response code
}

// FauxResponse represents a fake response to return
type FauxResponse struct {
	StatusCode int               // HTTP status code
	Headers    map[string]string // Response headers
	Body       string            // Response body
}

// GenerateFauxRequest generates a fake request for a given call flow step
func (f *FauxCallGenerator) GenerateFauxRequest(flowID string, stepNum int, componentName string, baseURL string) (*FauxRequest, error) {
	flow, exists := f.registry.GetFlow(flowID)
	if !exists {
		return nil, fmt.Errorf("call flow %s not found", flowID)
	}

	var step *Step
	for i := range flow.Steps {
		if flow.Steps[i].Sequence == stepNum {
			step = &flow.Steps[i]
			break
		}
	}

	if step == nil {
		return nil, fmt.Errorf("step %d not found in flow %s", stepNum, flowID)
	}

	// Determine if this step is for our component
	compSteps := f.registry.GetComponentSteps(componentName, flowID)
	isOurStep := false
	for _, s := range compSteps {
		if s.Sequence == stepNum {
			isOurStep = true
			break
		}
	}

	// Generate request based on step
	req := &FauxRequest{
		Method: "POST", // SIP messages are typically POST
		Headers: map[string]string{
			"Content-Type": "application/sip",
		},
		ExpectedCode: 200,
	}

	// Generate SIP message body based on step
	switch step.Message {
	case "REGISTER":
		req.Body = f.generateSIPRegister(step, componentName)
	case "INVITE":
		req.Body = f.generateSIPInvite(step, componentName)
	case "200 OK":
		req.Body = f.generateSIP200OK(step, componentName)
		req.ExpectedCode = 200
	case "180 Ringing":
		req.Body = f.generateSIP180Ringing(step, componentName)
		req.ExpectedCode = 180
	case "401 Unauthorized":
		req.Body = f.generateSIP401Unauthorized(step, componentName)
		req.ExpectedCode = 401
	default:
		req.Body = f.generateGenericSIPMessage(step, componentName)
	}

	// Determine target URL based on step direction and component
	if step.Direction == "request" && isOurStep {
		// This is a request we're sending out
		req.URL = f.getNeighborURL(componentName, step.To, baseURL)
	} else if step.Direction == "response" && isOurStep {
		// This is a response we're sending back
		req.URL = baseURL // Response goes back to caller
	} else {
		// This is a request coming to us
		req.URL = baseURL
	}

	return req, nil
}

// GenerateFauxResponse generates a fake response for a given call flow step
func (f *FauxCallGenerator) GenerateFauxResponse(flowID string, stepNum int, componentName string) (*FauxResponse, error) {
	flow, exists := f.registry.GetFlow(flowID)
	if !exists {
		return nil, fmt.Errorf("call flow %s not found", flowID)
	}

	var step *Step
	for i := range flow.Steps {
		if flow.Steps[i].Sequence == stepNum {
			step = &flow.Steps[i]
			break
		}
	}

	if step == nil {
		return nil, fmt.Errorf("step %d not found in flow %s", stepNum, flowID)
	}

	resp := &FauxResponse{
		Headers: map[string]string{
			"Content-Type": "application/sip",
		},
	}

	// Generate response based on step
	switch step.Message {
	case "200 OK":
		resp.StatusCode = 200
		resp.Body = f.generateSIP200OK(step, componentName)
	case "180 Ringing":
		resp.StatusCode = 180
		resp.Body = f.generateSIP180Ringing(step, componentName)
	case "401 Unauthorized":
		resp.StatusCode = 401
		resp.Body = f.generateSIP401Unauthorized(step, componentName)
	case "UAA", "SAA":
		// HSS responses
		resp.StatusCode = 200
		resp.Body = f.generateHSSResponse(step, componentName)
	default:
		resp.StatusCode = 200
		resp.Body = f.generateGenericSIPMessage(step, componentName)
	}

	return resp, nil
}

// ExecuteFauxRequest executes a fake request using curl-like HTTP call
func (f *FauxCallGenerator) ExecuteFauxRequest(req *FauxRequest) (*FauxResponse, error) {
	httpReq, err := http.NewRequest(req.Method, req.URL, bytes.NewBufferString(req.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Execute request
	resp, err := f.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Convert response headers
	respHeaders := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			respHeaders[k] = v[0]
		}
	}

	return &FauxResponse{
		StatusCode: resp.StatusCode,
		Headers:    respHeaders,
		Body:       string(body),
	}, nil
}

// Helper methods for generating SIP messages

func (f *FauxCallGenerator) generateSIPRegister(step *Step, componentName string) string {
	return fmt.Sprintf(`REGISTER sip:%s SIP/2.0
Via: SIP/2.0/UDP %s;branch=z9hG4bK-faux
From: <sip:test@example.com>;tag=faux-from
To: <sip:test@example.com>
Call-ID: faux-call-id-%d
CSeq: 1 REGISTER
Contact: <sip:test@example.com>
Content-Length: 0

`, componentName, step.From, step.Sequence)
}

func (f *FauxCallGenerator) generateSIPInvite(step *Step, componentName string) string {
	return fmt.Sprintf(`INVITE sip:test@example.com SIP/2.0
Via: SIP/2.0/UDP %s;branch=z9hG4bK-faux
From: <sip:caller@example.com>;tag=faux-from
To: <sip:callee@example.com>
Call-ID: faux-invite-%d
CSeq: 1 INVITE
Contact: <sip:caller@example.com>
Content-Type: application/sdp
Content-Length: 142

v=0
o=caller 2890844526 2890844526 IN IP4 192.0.2.1
s=-
c=IN IP4 192.0.2.1
t=0 0
m=audio 49170 RTP/AVP 0

`, step.From, step.Sequence)
}

func (f *FauxCallGenerator) generateSIP200OK(step *Step, componentName string) string {
	return fmt.Sprintf(`SIP/2.0 200 OK
Via: SIP/2.0/UDP %s;branch=z9hG4bK-faux
From: <sip:caller@example.com>;tag=faux-from
To: <sip:callee@example.com>;tag=faux-to
Call-ID: faux-call-id-%d
CSeq: 1 %s
Contact: <sip:callee@example.com>
Content-Type: application/sdp
Content-Length: 142

v=0
o=callee 2890844527 2890844527 IN IP4 192.0.2.2
s=-
c=IN IP4 192.0.2.2
t=0 0
m=audio 49172 RTP/AVP 0

`, step.From, step.Sequence, step.Message)
}

func (f *FauxCallGenerator) generateSIP180Ringing(step *Step, componentName string) string {
	return fmt.Sprintf(`SIP/2.0 180 Ringing
Via: SIP/2.0/UDP %s;branch=z9hG4bK-faux
From: <sip:caller@example.com>;tag=faux-from
To: <sip:callee@example.com>;tag=faux-to
Call-ID: faux-call-id-%d
CSeq: 1 INVITE
Content-Length: 0

`, step.From, step.Sequence)
}

func (f *FauxCallGenerator) generateSIP401Unauthorized(step *Step, componentName string) string {
	return fmt.Sprintf(`SIP/2.0 401 Unauthorized
Via: SIP/2.0/UDP %s;branch=z9hG4bK-faux
From: <sip:test@example.com>;tag=faux-from
To: <sip:test@example.com>;tag=faux-to
Call-ID: faux-call-id-%d
CSeq: 1 REGISTER
WWW-Authenticate: Digest realm="example.com", nonce="faux-nonce", algorithm=AKAv1-MD5
Content-Length: 0

`, step.From, step.Sequence)
}

func (f *FauxCallGenerator) generateHSSResponse(step *Step, componentName string) string {
	// HSS responses are typically Diameter/Cx protocol, but we'll simulate with JSON
	response := map[string]interface{}{
		"result_code": 2001, // DIAMETER_SUCCESS
		"scscf":       "scscf.example.com",
		"service_profile": map[string]interface{}{
			"public_identity": "sip:test@example.com",
			"ifc": []map[string]interface{}{
				{
					"priority": 1,
					"trigger_point": "INVITE",
					"application_server": "as.example.com",
				},
			},
		},
	}

	jsonData, _ := json.Marshal(response)
	return string(jsonData)
}

func (f *FauxCallGenerator) generateGenericSIPMessage(step *Step, componentName string) string {
	return fmt.Sprintf(`%s sip:test@example.com SIP/2.0
Via: SIP/2.0/UDP %s;branch=z9hG4bK-faux
From: <sip:test@example.com>;tag=faux-from
To: <sip:test@example.com>
Call-ID: faux-call-id-%d
CSeq: 1 %s
Content-Length: 0

`, step.Message, step.From, step.Sequence, step.Message)
}

func (f *FauxCallGenerator) getNeighborURL(componentName, neighbor, baseURL string) string {
	// Map neighbors to typical ports/components
	neighborPorts := map[string]string{
		"UE":      "http://localhost:5060",
		"P-CSCF":  "http://localhost:8081",
		"I-CSCF":  "http://localhost:8082",
		"S-CSCF":  "http://localhost:8083",
		"BGCF":    "http://localhost:8084",
		"MGCF":    "http://localhost:8085",
		"HSS":     "http://localhost:8086",
		"PSTN":    "http://localhost:5061",
		"Destination": "http://localhost:5062",
	}

	if url, exists := neighborPorts[neighbor]; exists {
		return url
	}

	// Default to baseURL if neighbor not found
	return baseURL
}
