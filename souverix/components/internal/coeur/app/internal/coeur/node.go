package coeur

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Node represents the base interface for all IMS core nodes
// All IMS nodes (P-CSCF, I-CSCF, S-CSCF, BGCF, MGCF, MGW) implement this interface
type Node interface {
	// Start starts the node
	Start() error

	// Stop stops the node gracefully
	Stop() error

	// Name returns the node name (e.g., "pcscf", "icscf", "scscf")
	Name() string

	// Status returns the current node status
	Status() NodeStatus
}

// NodeStatus represents the operational status of a node
type NodeStatus struct {
	State      NodeState `json:"state"`
	StartedAt  string    `json:"started_at,omitempty"`
	Uptime     string    `json:"uptime,omitempty"`
	LastError  string    `json:"last_error,omitempty"`
	Metrics    NodeMetrics `json:"metrics,omitempty"`
}

// NodeState represents the state of a node
type NodeState string

const (
	NodeStateStopped  NodeState = "stopped"
	NodeStateStarting NodeState = "starting"
	NodeStateRunning  NodeState = "running"
	NodeStateStopping NodeState = "stopping"
	NodeStateError    NodeState = "error"
)

// NodeMetrics contains basic metrics for a node
type NodeMetrics struct {
	MessagesProcessed uint64 `json:"messages_processed"`
	MessagesFailed    uint64 `json:"messages_failed"`
	ActiveSessions    uint64 `json:"active_sessions"`
	LastMessageAt     string `json:"last_message_at,omitempty"`
}

// BaseNode provides common functionality for all IMS nodes
type BaseNode struct {
	name      string
	state     NodeState
	startedAt time.Time
	metrics   NodeMetrics
	mu        sync.RWMutex
	log       *logrus.Logger
}

// NewBaseNode creates a new base node
func NewBaseNode(name string, log *logrus.Logger) *BaseNode {
	return &BaseNode{
		name:  name,
		state: NodeStateStopped,
		log:   log,
	}
}

// Name returns the node name
func (n *BaseNode) Name() string {
	return n.name
}

// Status returns the current node status
func (n *BaseNode) Status() NodeStatus {
	n.mu.RLock()
	defer n.mu.RUnlock()

	status := NodeStatus{
		State:   n.state,
		Metrics: n.metrics,
	}

	if !n.startedAt.IsZero() {
		status.StartedAt = n.startedAt.Format(time.RFC3339)
		status.Uptime = time.Since(n.startedAt).String()
	}

	return status
}

// setState updates the node state
func (n *BaseNode) setState(state NodeState) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.state = state
}

// recordMessage records a processed message
func (n *BaseNode) recordMessage(success bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if success {
		n.metrics.MessagesProcessed++
	} else {
		n.metrics.MessagesFailed++
	}
	n.metrics.LastMessageAt = time.Now().Format(time.RFC3339)
}

// recordSession records an active session change
func (n *BaseNode) recordSession(delta int64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if delta > 0 {
		n.metrics.ActiveSessions += uint64(delta)
	} else {
		// Prevent underflow
		if uint64(-delta) <= n.metrics.ActiveSessions {
			n.metrics.ActiveSessions -= uint64(-delta)
		}
	}
}
