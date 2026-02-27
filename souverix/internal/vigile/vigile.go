package vigile

import (
	"context"
	
	"github.com/dasmlab/ims/internal/gouverne/config"
	"github.com/sirupsen/logrus"
)

// Component represents the Souverix Vigile component
type Component struct {
	log *logrus.Logger
}

// New creates a new component instance
func New(cfg *config.Config, log *logrus.Logger) *Component {
	return &Component{
		log: log,
	}
}

// Start starts the component
func (c *Component) Start(ctx context.Context) error {
	c.log.Info("vigile component started (stub)")
	return nil
}

// Stop stops the component
func (c *Component) Stop(ctx context.Context) error {
	c.log.Info("vigile component stopped (stub)")
	return nil
}
