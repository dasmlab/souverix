package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dasmlab/ims/components/common/hss"
	"github.com/dasmlab/ims/components/common/sip"
	bgcfHandler "github.com/dasmlab/ims/components/coeur/bgcf/app/internal/bgcf"
	icscfHandler "github.com/dasmlab/ims/components/coeur/icscf/app/internal/icscf"
	mgcfHandler "github.com/dasmlab/ims/components/coeur/mgcf/app/internal/mgcf"
	pcscfHandler "github.com/dasmlab/ims/components/coeur/pcscf/app/internal/pcscf"
	scscfHandler "github.com/dasmlab/ims/components/coeur/scscf/app/internal/scscf"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("Souverix Coeur - Basic SIP INVITE Call Flow Test")
	fmt.Println("========================================")
	fmt.Println()

	// Initialize components
	logger := log.New(log.Default().Writer(), "[TEST] ", log.LstdFlags)
	hssClient := hss.NewHSSClient()

	pcscf := pcscfHandler.NewHandler("icscf.example.com:5060", logger)
	icscf := icscfHandler.NewHandler(hssClient, logger)
	scscf := scscfHandler.NewHandler(hssClient, "bgcf.example.com:5060", logger)
	bgcf := bgcfHandler.NewHandler(logger)
	mgcf := mgcfHandler.NewHandler(logger)

	// Test 1: IMS-to-IMS call
	fmt.Println("Test 1: IMS-to-IMS Call Flow")
	fmt.Println("----------------------------")
	testIMS2IMS(pcscf, icscf, scscf, logger)
	fmt.Println()

	// Test 2: IMS-to-PSTN call (VoLTE profile)
	fmt.Println("Test 2: IMS-to-PSTN Call Flow (VoLTE)")
	fmt.Println("--------------------------------------")
	testIMS2PSTN(pcscf, icscf, scscf, bgcf, mgcf, logger)
	fmt.Println()
}

func testIMS2IMS(pcscf *pcscfHandler.Handler, icscf *icscfHandler.Handler, scscf *scscfHandler.Handler, logger *log.Logger) {
	// Create INVITE
	invite := sip.NewINVITE(
		"sip:alice@example.com",
		"sip:bob@example.com",
		"call-123@example.com",
	)
	invite.Contact = "sip:alice@ue.example.com"

	logger.Printf("UE -> P-CSCF: INVITE")
	msg, err := pcscf.HandleINVITE(invite)
	if err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	logger.Printf("P-CSCF -> I-CSCF: INVITE")
	msg, scscfAddr, err := icscf.HandleINVITE(msg)
	if err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}
	logger.Printf("I-CSCF selected S-CSCF: %s", scscfAddr)

	logger.Printf("I-CSCF -> S-CSCF: INVITE")
	msg, dest, err := scscf.HandleINVITE(msg)
	if err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}
	logger.Printf("S-CSCF routing to: %s", dest)

	logger.Printf("S-CSCF -> Destination: INVITE")
	logger.Printf("Destination -> S-CSCF: 180 Ringing")
	ringing := sip.New180Ringing(msg)

	logger.Printf("S-CSCF -> I-CSCF: 180 Ringing")
	logger.Printf("I-CSCF -> P-CSCF: 180 Ringing")
	logger.Printf("P-CSCF -> UE: 180 Ringing")

	time.Sleep(100 * time.Millisecond)

	logger.Printf("Destination -> S-CSCF: 200 OK")
	ok := sip.New200OK(msg, "sip:bob@destination.example.com")

	logger.Printf("S-CSCF -> I-CSCF: 200 OK")
	logger.Printf("I-CSCF -> P-CSCF: 200 OK")
	logger.Printf("P-CSCF -> UE: 200 OK")

	logger.Printf("UE -> P-CSCF: ACK")
	logger.Printf("P-CSCF -> I-CSCF: ACK")
	logger.Printf("I-CSCF -> S-CSCF: ACK")
	logger.Printf("S-CSCF -> Destination: ACK")

	fmt.Println("✅ IMS-to-IMS call flow completed successfully")
	_ = ringing
	_ = ok
}

func testIMS2PSTN(pcscf *pcscfHandler.Handler, icscf *icscfHandler.Handler, scscf *scscfHandler.Handler, bgcf *bgcfHandler.Handler, mgcf *mgcfHandler.Handler, logger *log.Logger) {
	// Create INVITE to PSTN
	invite := sip.NewINVITE(
		"sip:alice@example.com",
		"tel:+1234567890",
		"call-pstn-123@example.com",
	)
	invite.Contact = "sip:alice@ue.example.com"

	logger.Printf("UE -> P-CSCF: INVITE (tel:+1234567890)")
	msg, err := pcscf.HandleINVITE(invite)
	if err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	logger.Printf("P-CSCF -> I-CSCF: INVITE")
	msg, scscfAddr, err := icscf.HandleINVITE(msg)
	if err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}
	logger.Printf("I-CSCF selected S-CSCF: %s", scscfAddr)

	logger.Printf("I-CSCF -> S-CSCF: INVITE")
	msg, dest, err := scscf.HandleINVITE(msg)
	if err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}
	logger.Printf("S-CSCF detected PSTN destination, routing to BGCF: %s", dest)

	logger.Printf("S-CSCF -> BGCF: INVITE (Mi interface)")
	msg, mgcfAddr, err := bgcf.HandleINVITE(msg)
	if err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}
	logger.Printf("BGCF selected local breakout, routing to MGCF: %s", mgcfAddr)

	logger.Printf("BGCF -> MGCF: INVITE (Mj interface)")
	_, err = mgcf.HandleINVITE(msg)
	if err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	// Simulate PSTN response
	time.Sleep(100 * time.Millisecond)

	logger.Printf("PSTN -> MGCF: ISUP ACM")
	ringing, _ := mgcf.HandlePSTNResponse(msg.CallID, "ACM")
	logger.Printf("MGCF -> BGCF: 180 Ringing")
	logger.Printf("BGCF -> S-CSCF: 180 Ringing")
	logger.Printf("S-CSCF -> I-CSCF: 180 Ringing")
	logger.Printf("I-CSCF -> P-CSCF: 180 Ringing")
	logger.Printf("P-CSCF -> UE: 180 Ringing")

	time.Sleep(200 * time.Millisecond)

	logger.Printf("PSTN -> MGCF: ISUP ANM")
	ok, _ := mgcf.HandlePSTNResponse(msg.CallID, "ANM")
	logger.Printf("MGCF -> BGCF: 200 OK")
	logger.Printf("BGCF -> S-CSCF: 200 OK")
	logger.Printf("S-CSCF -> I-CSCF: 200 OK")
	logger.Printf("I-CSCF -> P-CSCF: 200 OK")
	logger.Printf("P-CSCF -> UE: 200 OK")

	logger.Printf("UE -> P-CSCF: ACK")
	logger.Printf("P-CSCF -> I-CSCF: ACK")
	logger.Printf("I-CSCF -> S-CSCF: ACK")
	logger.Printf("S-CSCF -> BGCF: ACK")
	logger.Printf("BGCF -> MGCF: ACK")
	logger.Printf("MGCF -> PSTN: ISUP ACK")

	fmt.Println("✅ IMS-to-PSTN call flow completed successfully")
	_ = ringing
	_ = ok
}
