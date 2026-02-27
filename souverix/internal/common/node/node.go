// Package node provides the base interface and boilerplate for all IMS core nodes.
package node

import (
	"context"
	"time"
)

// Node represents an IMS core node (P-CSCF, I-CSCF, S-CSCF, BGCF, MGCF, HSS, MGW).
type Node interface {
	// Name returns the node name (e.g., "pcscf", "scscf").
	Name() string

	// Start initializes and starts the node.
	Start(ctx context.Context) error

	// Stop gracefully stops the node.
	Stop(ctx context.Context) error

	// Health returns the current health status.
	Health() HealthStatus

	// Metrics returns node-specific metrics.
	Metrics() Metrics
}

// HealthStatus represents the health of a node.
type HealthStatus struct {
	Status    string    // "healthy", "degraded", "unhealthy"
	Timestamp time.Time
	Details   map[string]interface{}
}

// Metrics represents node-specific metrics.
type Metrics struct {
	MessagesProcessed uint64
	Errors            uint64
	LatencyP50        time.Duration
	LatencyP95        time.Duration
	LatencyP99        time.Duration
}

// BaseNode provides common functionality for all nodes.
type BaseNode struct {
	name      string
	startedAt time.Time
	health    HealthStatus
	metrics   Metrics
}

// NewBaseNode creates a new base node.
func NewBaseNode(name string) *BaseNode {
	return &BaseNode{
		name:      name,
		health:    HealthStatus{Status: "unhealthy", Timestamp: time.Now()},
		metrics:   Metrics{},
	}
}

// Name returns the node name.
func (b *BaseNode) Name() string {
	return b.name
}

// Health returns the current health status.
func (b *BaseNode) Health() HealthStatus {
	return b.health
}

// Metrics returns node-specific metrics.
func (b *BaseNode) Metrics() Metrics {
	return b.metrics
}

// SetHealth updates the health status.
func (b *BaseNode) SetHealth(status string, details map[string]interface{}) {
	b.health = HealthStatus{
		Status:    status,
		Timestamp: time.Now(),
		Details:   details,
	}
}

// IncrementMessages increments the messages processed counter.
func (b *BaseNode) IncrementMessages() {
	b.metrics.MessagesProcessed++
}

// IncrementErrors increments the error counter.
func (b *BaseNode) IncrementErrors() {
	b.metrics.Errors++
}
