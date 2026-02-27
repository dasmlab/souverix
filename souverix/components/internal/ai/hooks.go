package ai

import (
	"github.com/dasmlab/ims/internal/sip"
	"github.com/sirupsen/logrus"
)

// AgentHook represents a hook point for AI agent integration
type AgentHook interface {
	// OnSIPMessage is called when a SIP message is received
	OnSIPMessage(msg *sip.Message) error

	// OnSessionStart is called when a new session starts
	OnSessionStart(sessionID string) error

	// OnSessionEnd is called when a session ends
	OnSessionEnd(sessionID string) error

	// OnRegistration is called when a registration occurs
	OnRegistration(impi, impu string) error
}

// MCPAgent represents an MCP (Model Context Protocol) agent
type MCPAgent struct {
	name   string
	log    *logrus.Logger
	client MCPClient // TODO: Implement MCP client
}

// MCPClient represents an MCP client interface
type MCPClient interface {
	SendMessage(msg interface{}) error
	ReceiveMessage() (interface{}, error)
}

// NewMCPAgent creates a new MCP agent hook
func NewMCPAgent(name string, log *logrus.Logger) *MCPAgent {
	return &MCPAgent{
		name: name,
		log:  log,
	}
}

// OnSIPMessage handles SIP message events
func (a *MCPAgent) OnSIPMessage(msg *sip.Message) error {
	a.log.WithFields(logrus.Fields{
		"agent":  a.name,
		"method": msg.Method,
	}).Debug("AI agent processing SIP message")

	// TODO: Send to MCP client for processing
	// This is a placeholder for future MCP integration
	return nil
}

// OnSessionStart handles session start events
func (a *MCPAgent) OnSessionStart(sessionID string) error {
	a.log.WithFields(logrus.Fields{
		"agent":     a.name,
		"sessionID": sessionID,
	}).Debug("AI agent notified of session start")

	// TODO: Send to MCP client
	return nil
}

// OnSessionEnd handles session end events
func (a *MCPAgent) OnSessionEnd(sessionID string) error {
	a.log.WithFields(logrus.Fields{
		"agent":     a.name,
		"sessionID": sessionID,
	}).Debug("AI agent notified of session end")

	// TODO: Send to MCP client
	return nil
}

// OnRegistration handles registration events
func (a *MCPAgent) OnRegistration(impi, impu string) error {
	a.log.WithFields(logrus.Fields{
		"agent": a.name,
		"impi":  impi,
		"impu":  impu,
	}).Debug("AI agent notified of registration")

	// TODO: Send to MCP client
	return nil
}

// HookManager manages AI agent hooks
type HookManager struct {
	hooks []AgentHook
	log   *logrus.Logger
}

// NewHookManager creates a new hook manager
func NewHookManager(log *logrus.Logger) *HookManager {
	return &HookManager{
		hooks: make([]AgentHook, 0),
		log:   log,
	}
}

// RegisterHook registers an AI agent hook
func (hm *HookManager) RegisterHook(hook AgentHook) {
	hm.hooks = append(hm.hooks, hook)
	hm.log.Info("AI agent hook registered")
}

// NotifySIPMessage notifies all hooks of a SIP message
func (hm *HookManager) NotifySIPMessage(msg *sip.Message) {
	for _, hook := range hm.hooks {
		if err := hook.OnSIPMessage(msg); err != nil {
			hm.log.WithError(err).Warn("AI agent hook error")
		}
	}
}

// NotifySessionStart notifies all hooks of a session start
func (hm *HookManager) NotifySessionStart(sessionID string) {
	for _, hook := range hm.hooks {
		if err := hook.OnSessionStart(sessionID); err != nil {
			hm.log.WithError(err).Warn("AI agent hook error")
		}
	}
}

// NotifySessionEnd notifies all hooks of a session end
func (hm *HookManager) NotifySessionEnd(sessionID string) {
	for _, hook := range hm.hooks {
		if err := hook.OnSessionEnd(sessionID); err != nil {
			hm.log.WithError(err).Warn("AI agent hook error")
		}
	}
}

// NotifyRegistration notifies all hooks of a registration
func (hm *HookManager) NotifyRegistration(impi, impu string) {
	for _, hook := range hm.hooks {
		if err := hook.OnRegistration(impi, impu); err != nil {
			hm.log.WithError(err).Warn("AI agent hook error")
		}
	}
}

// ExtensibilityPoint represents a point where AI agents can extend functionality
type ExtensibilityPoint string

const (
	ExtPointSIPMessage   ExtensibilityPoint = "sip_message"
	ExtPointSessionStart ExtensibilityPoint = "session_start"
	ExtPointSessionEnd   ExtensibilityPoint = "session_end"
	ExtPointRegistration ExtensibilityPoint = "registration"
	ExtPointCallRouting  ExtensibilityPoint = "call_routing"
	ExtPointFraudDetection ExtensibilityPoint = "fraud_detection"
)

// Extension represents an extension function
type Extension func(context interface{}) (interface{}, error)

// ExtensionRegistry manages extensions
type ExtensionRegistry struct {
	extensions map[ExtensibilityPoint][]Extension
	log        *logrus.Logger
}

// NewExtensionRegistry creates a new extension registry
func NewExtensionRegistry(log *logrus.Logger) *ExtensionRegistry {
	return &ExtensionRegistry{
		extensions: make(map[ExtensibilityPoint][]Extension),
		log:        log,
	}
}

// RegisterExtension registers an extension at a specific point
func (er *ExtensionRegistry) RegisterExtension(point ExtensibilityPoint, ext Extension) {
	er.extensions[point] = append(er.extensions[point], ext)
	er.log.WithField("point", point).Info("extension registered")
}

// ExecuteExtensions executes all extensions at a point
func (er *ExtensionRegistry) ExecuteExtensions(point ExtensibilityPoint, context interface{}) ([]interface{}, error) {
	extensions, ok := er.extensions[point]
	if !ok {
		return nil, nil
	}

	results := make([]interface{}, 0, len(extensions))
	for _, ext := range extensions {
		result, err := ext(context)
		if err != nil {
			er.log.WithError(err).WithField("point", point).Warn("extension execution error")
			continue
		}
		results = append(results, result)
	}

	return results, nil
}
