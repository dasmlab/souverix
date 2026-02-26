package li

import (
	"fmt"
	"sync"
	"time"

	"github.com/dasmlab/ims/internal/sip"
	"github.com/sirupsen/logrus"
)

// InterceptTarget represents a target for lawful interception
type InterceptTarget struct {
	IMPI        string
	IMPU        string
	TN          string
	WarrantID   string
	WarrantType string // "signaling", "signaling+media"
	Activated   bool
	ActivatedAt time.Time
	ExpiresAt   time.Time
}

// InterceptController manages lawful interception
type InterceptController struct {
	targets map[string]*InterceptTarget // key: IMPI or TN
	md      MediationDevice
	log     *logrus.Logger
	mu      sync.RWMutex
}

// MediationDevice represents the LI Mediation Device (MD)
type MediationDevice interface {
	// SendSignaling sends SIP signaling to MD
	SendSignaling(msg *sip.Message, targetID string) error

	// SendMedia sends RTP media to MD (if media intercept enabled)
	SendMedia(rtpData []byte, targetID string) error

	// IsAvailable checks if MD is available
	IsAvailable() bool
}

// NewInterceptController creates a new LI controller
func NewInterceptController(md MediationDevice, log *logrus.Logger) *InterceptController {
	return &InterceptController{
		targets: make(map[string]*InterceptTarget),
		md:      md,
		log:     log,
	}
}

// ActivateWarrant activates a lawful intercept warrant
func (ic *InterceptController) ActivateWarrant(target *InterceptTarget) error {
	ic.mu.Lock()
	defer ic.mu.Unlock()

	target.Activated = true
	target.ActivatedAt = time.Now()

	// Store by multiple keys for lookup
	ic.targets[target.IMPI] = target
	if target.IMPU != "" {
		ic.targets[target.IMPU] = target
	}
	if target.TN != "" {
		ic.targets[target.TN] = target
	}

	ic.log.WithFields(logrus.Fields{
		"warrant_id": target.WarrantID,
		"target":     target.IMPI,
		"type":       target.WarrantType,
	}).Info("LI warrant activated")

	return nil
}

// DeactivateWarrant deactivates a lawful intercept warrant
func (ic *InterceptController) DeactivateWarrant(warrantID string) error {
	ic.mu.Lock()
	defer ic.mu.Unlock()

	for key, target := range ic.targets {
		if target.WarrantID == warrantID {
			target.Activated = false
			delete(ic.targets, key)

			ic.log.WithFields(logrus.Fields{
				"warrant_id": warrantID,
				"target":     target.IMPI,
			}).Info("LI warrant deactivated")

			return nil
		}
	}

	return fmt.Errorf("warrant not found: %s", warrantID)
}

// IsTarget checks if a subscriber is under interception
func (ic *InterceptController) IsTarget(identifier string) (*InterceptTarget, bool) {
	ic.mu.RLock()
	defer ic.mu.RUnlock()

	target, ok := ic.targets[identifier]
	if !ok || !target.Activated {
		return nil, false
	}

	// Check if warrant expired
	if !target.ExpiresAt.IsZero() && time.Now().After(target.ExpiresAt) {
		return nil, false
	}

	return target, true
}

// InterceptMessage intercepts a SIP message if target is under warrant
func (ic *InterceptController) InterceptMessage(msg *sip.Message, fromID, toID string) error {
	// Check if either party is under interception
	var target *InterceptTarget
	var isTarget bool

	if target, isTarget = ic.IsTarget(fromID); !isTarget {
		if target, isTarget = ic.IsTarget(toID); !isTarget {
			return nil // Not a target, no interception
		}
	}

	if !ic.md.IsAvailable() {
		ic.log.Warn("LI mediation device unavailable")
		return fmt.Errorf("LI mediation device unavailable")
	}

	// Intercept signaling
	if err := ic.md.SendSignaling(msg, target.WarrantID); err != nil {
		ic.log.WithError(err).Error("failed to send signaling to LI MD")
		return err
	}

	ic.log.WithFields(logrus.Fields{
		"warrant_id": target.WarrantID,
		"call_id":    msg.GetHeader("Call-ID"),
		"method":     msg.Method,
	}).Debug("LI signaling intercepted")

	return nil
}

// InterceptMedia intercepts RTP media if target is under warrant
func (ic *InterceptController) InterceptMedia(rtpData []byte, targetID string) error {
	target, isTarget := ic.IsTarget(targetID)
	if !isTarget {
		return nil // Not a target
	}

	// Only intercept media if warrant type includes media
	if target.WarrantType != "signaling+media" {
		return nil
	}

	if !ic.md.IsAvailable() {
		return fmt.Errorf("LI mediation device unavailable")
	}

	return ic.md.SendMedia(rtpData, target.WarrantID)
}

// ListActiveWarrants returns all active warrants
func (ic *InterceptController) ListActiveWarrants() []*InterceptTarget {
	ic.mu.RLock()
	defer ic.mu.RUnlock()

	active := make([]*InterceptTarget, 0)
	for _, target := range ic.targets {
		if target.Activated {
			active = append(active, target)
		}
	}

	return active
}

// AuditLog represents an audit log entry for LI
type AuditLog struct {
	Timestamp  time.Time
	WarrantID  string
	Action     string // "activate", "deactivate", "intercept"
	Target     string
	CallID     string
	Operator   string
	Details    string
}

// AuditLogger logs LI operations for compliance
type AuditLogger struct {
	logs []AuditLog
	mu   sync.RWMutex
	log  *logrus.Logger
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(log *logrus.Logger) *AuditLogger {
	return &AuditLogger{
		logs: make([]AuditLog, 0),
		log:  log,
	}
}

// Log logs an audit event
func (al *AuditLogger) Log(entry AuditLog) {
	al.mu.Lock()
	defer al.mu.Unlock()

	entry.Timestamp = time.Now()
	al.logs = append(al.logs, entry)

	// Log to structured logger
	al.log.WithFields(logrus.Fields{
		"warrant_id": entry.WarrantID,
		"action":     entry.Action,
		"target":     entry.Target,
		"call_id":    entry.CallID,
		"operator":   entry.Operator,
	}).Info("LI audit log")
}

// GetLogs returns audit logs (with access control in production)
func (al *AuditLogger) GetLogs(warrantID string, startTime, endTime time.Time) []AuditLog {
	al.mu.RLock()
	defer al.mu.RUnlock()

	filtered := make([]AuditLog, 0)
	for _, log := range al.logs {
		if log.WarrantID == warrantID &&
			log.Timestamp.After(startTime) &&
			log.Timestamp.Before(endTime) {
			filtered = append(filtered, log)
		}
	}

	return filtered
}
