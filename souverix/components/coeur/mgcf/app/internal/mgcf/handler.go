package mgcf

import (
	"fmt"
	"log"
	
	"github.com/dasmlab/souverix/common/sip"
)

// Handler handles SIP messages in MGCF
type Handler struct {
	logger *log.Logger
}

// NewHandler creates a new MGCF handler
func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// HandleINVITE processes an INVITE request for PSTN interworking
func (h *Handler) HandleINVITE(msg *sip.Message) (*sip.Message, error) {
	h.logger.Printf("MGCF: Received INVITE for PSTN interworking from %s to %s", msg.From, msg.To)
	
	// Convert SIP INVITE to ISUP IAM
	h.logger.Printf("MGCF: Converting SIP INVITE to ISUP IAM")
	
	// Control MGW via H.248 (Mg interface)
	h.logger.Printf("MGCF: Sending H.248 Add command to MGW")
	
	// Send ISUP IAM to PSTN (Nc interface)
	h.logger.Printf("MGCF: Sending ISUP IAM to PSTN")
	
	// For now, simulate PSTN response
	// In real implementation, wait for ISUP ACM/ANM
	return nil, nil
}

// HandlePSTNResponse simulates PSTN response (ISUP ACM/ANM)
func (h *Handler) HandlePSTNResponse(callID string, responseType string) (*sip.Message, error) {
	h.logger.Printf("MGCF: Received PSTN response: %s for Call-ID: %s", responseType, callID)
	
	switch responseType {
	case "ACM":
		// Convert ISUP ACM to SIP 180 Ringing
		h.logger.Printf("MGCF: Converting ISUP ACM to SIP 180 Ringing")
		return &sip.Message{
			Method: "180",
			CallID: callID,
		}, nil
	case "ANM":
		// Convert ISUP ANM to SIP 200 OK
		h.logger.Printf("MGCF: Converting ISUP ANM to SIP 200 OK")
		return &sip.Message{
			Method: "200",
			CallID: callID,
		}, nil
	default:
		return nil, fmt.Errorf("unknown PSTN response type: %s", responseType)
	}
}
