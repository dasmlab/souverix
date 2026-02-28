package pcscf

import (
	"fmt"
	"log"
	
	"github.com/dasmlab/souverix/common/sip"
)

// Handler handles SIP messages in P-CSCF
type Handler struct {
	nextHop string // I-CSCF address
	logger  *log.Logger
}

// NewHandler creates a new P-CSCF handler
func NewHandler(icscfAddress string, logger *log.Logger) *Handler {
	return &Handler{
		nextHop: icscfAddress,
		logger:  logger,
	}
}

// HandleINVITE processes an INVITE request
func (h *Handler) HandleINVITE(msg *sip.Message) (*sip.Message, error) {
	h.logger.Printf("P-CSCF: Received INVITE from %s to %s", msg.From, msg.To)
	
	// Validate SIP headers
	if err := h.validateHeaders(msg); err != nil {
		return nil, fmt.Errorf("header validation failed: %w", err)
	}
	
	// Insert Record-Route to remain in path
	msg.AddRecordRoute("<sip:pcscf.example.com;lr>")
	
	// Forward to I-CSCF (Mw interface)
	h.logger.Printf("P-CSCF: Forwarding INVITE to I-CSCF at %s", h.nextHop)
	
	// In real implementation, this would send over network
	// For now, we return the message with Record-Route added
	return msg, nil
}

// HandleResponse processes a SIP response
func (h *Handler) HandleResponse(msg *sip.Message) (*sip.Message, error) {
	h.logger.Printf("P-CSCF: Received %s response for Call-ID: %s", msg.Method, msg.CallID)
	
	// Forward response back to UE
	return msg, nil
}

func (h *Handler) validateHeaders(msg *sip.Message) error {
	if msg.From == "" {
		return fmt.Errorf("missing From header")
	}
	if msg.To == "" {
		return fmt.Errorf("missing To header")
	}
	if msg.CallID == "" {
		return fmt.Errorf("missing Call-ID header")
	}
	return nil
}
