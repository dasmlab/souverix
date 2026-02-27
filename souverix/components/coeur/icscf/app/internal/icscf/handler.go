package icscf

import (
	"fmt"
	"log"
	
	"github.com/dasmlab/ims/components/common/hss"
	"github.com/dasmlab/ims/components/common/sip"
)

// Handler handles SIP messages in I-CSCF
type Handler struct {
	hssClient *hss.HSSClient
	scscfPool []string
	logger    *log.Logger
}

// NewHandler creates a new I-CSCF handler
func NewHandler(hssClient *hss.HSSClient, logger *log.Logger) *Handler {
	return &Handler{
		hssClient: hssClient,
		scscfPool: []string{"scscf1.example.com", "scscf2.example.com"},
		logger:    logger,
	}
}

// HandleINVITE processes an INVITE request
func (h *Handler) HandleINVITE(msg *sip.Message) (*sip.Message, string, error) {
	h.logger.Printf("I-CSCF: Received INVITE from %s to %s", msg.From, msg.To)
	
	// Extract user identity from From header
	impi := h.extractIMPI(msg.From)
	
	// Query HSS for S-CSCF assignment (UAR/UAA)
	scscf, capabilities, err := h.hssClient.GetSCSCFAssignment(impi)
	if err != nil {
		h.logger.Printf("I-CSCF: HSS query failed, using default S-CSCF")
		scscf = h.scscfPool[0] // Fallback to first S-CSCF
	}
	
	h.logger.Printf("I-CSCF: HSS assigned S-CSCF: %s (capabilities: %v)", scscf, capabilities)
	
	// Forward to selected S-CSCF (Mw interface)
	h.logger.Printf("I-CSCF: Forwarding INVITE to S-CSCF at %s", scscf)
	
	return msg, scscf, nil
}

// HandleResponse processes a SIP response
func (h *Handler) HandleResponse(msg *sip.Message) (*sip.Message, error) {
	h.logger.Printf("I-CSCF: Received %s response for Call-ID: %s", msg.Method, msg.CallID)
	
	// Forward response back (topology hiding)
	return msg, nil
}

func (h *Handler) extractIMPI(from string) string {
	// Simplified extraction - in real implementation, parse SIP URI properly
	if from == "" {
		return "sip:user@example.com"
	}
	return from
}
