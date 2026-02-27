package scscf

import (
	"fmt"
	"log"
	"strings"
	
	"github.com/dasmlab/ims/components/common/hss"
	"github.com/dasmlab/ims/components/common/sip"
)

// Handler handles SIP messages in S-CSCF
type Handler struct {
	hssClient  *hss.HSSClient
	bgcfAddress string
	logger     *log.Logger
}

// NewHandler creates a new S-CSCF handler
func NewHandler(hssClient *hss.HSSClient, bgcfAddress string, logger *log.Logger) *Handler {
	return &Handler{
		hssClient:  hssClient,
		bgcfAddress: bgcfAddress,
		logger:     logger,
	}
}

// HandleINVITE processes an INVITE request
func (h *Handler) HandleINVITE(msg *sip.Message) (*sip.Message, string, error) {
	h.logger.Printf("S-CSCF: Received INVITE from %s to %s", msg.From, msg.To)
	
	// Extract user identity
	impi := h.extractIMPI(msg.From)
	
	// Get service profile from HSS (SAR/SAA)
	profile, err := h.hssClient.GetServiceProfile(impi)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get service profile: %w", err)
	}
	
	h.logger.Printf("S-CSCF: Service profile loaded for %s (registered: %v)", impi, profile.Registered)
	
	// Apply Initial Filter Criteria (iFC)
	h.logger.Printf("S-CSCF: Applying iFC (count: %d)", len(profile.IFC))
	for _, ifc := range profile.IFC {
		if ifc.Trigger == "INVITE" && ifc.ApplicationServer != "" {
			h.logger.Printf("S-CSCF: Triggering AS: %s", ifc.ApplicationServer)
			// In real implementation, forward to AS via ISC interface
		}
	}
	
	// Check for call barring
	if profile.Barring {
		return nil, "", fmt.Errorf("call barred for user: %s", impi)
	}
	
	// Insert Record-Route to anchor dialog
	msg.AddRecordRoute("<sip:scscf.example.com;lr>")
	
	// Determine routing
	destination := msg.To
	if msg.IsTelURI() {
		// PSTN breakout required - route to BGCF
		h.logger.Printf("S-CSCF: PSTN destination detected, routing to BGCF")
		return msg, h.bgcfAddress, nil
	}
	
	// IMS destination - route to terminating side
	h.logger.Printf("S-CSCF: IMS destination, routing to: %s", destination)
	return msg, destination, nil
}

// HandleResponse processes a SIP response
func (h *Handler) HandleResponse(msg *sip.Message) (*sip.Message, error) {
	h.logger.Printf("S-CSCF: Received %s response for Call-ID: %s", msg.Method, msg.CallID)
	
	// Forward response back through Record-Route
	return msg, nil
}

func (h *Handler) extractIMPI(from string) string {
	if from == "" {
		return "sip:user@example.com"
	}
	return from
}
