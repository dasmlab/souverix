package ims

// Subscriber represents an IMS subscriber
type Subscriber struct {
	// IMS Public User Identity (IMPU)
	IMPU string

	// IMS Private User Identity (IMPI)
	IMPI string

	// Registration state
	Registered bool

	// Contact information
	Contact string

	// S-CSCF assignment
	SCSCFName string

	// Service profile
	ServiceProfile ServiceProfile

	// Authentication data
	AuthData AuthData
}

// ServiceProfile represents a subscriber's service profile
type ServiceProfile struct {
	PublicIdentities []string
	CoreNetworkServices []string
	InitialFilterCriteria []FilterCriteria
}

// FilterCriteria represents Initial Filter Criteria (iFC) for AS routing
type FilterCriteria struct {
	Priority    int
	Trigger     TriggerPoint
	ApplicationServer ApplicationServer
}

// TriggerPoint represents a trigger point for AS invocation
type TriggerPoint struct {
	ConditionTypeCNF string
	SPT              []ServicePointTrigger
}

// ServicePointTrigger represents a service point trigger
type ServicePointTrigger struct {
	ConditionNegated bool
	Group            string
	Method           string
	RequestURI       string
}

// ApplicationServer represents an Application Server
type ApplicationServer struct {
	ServerName string
	DefaultHandling string // "SESSION_CONTINUED", "SESSION_TERMINATED"
}

// AuthData represents authentication data for a subscriber
type AuthData struct {
	AuthScheme string // "Digest", "AKA"
	Username   string
	Realm      string
	Password   string // Hashed
	HA1        string // Pre-computed HA1 for Digest
}

// Session represents an IMS session
type Session struct {
	SessionID string

	// Call-ID
	CallID string

	// From/To tags
	FromTag string
	ToTag   string

	// Request URI
	RequestURI string

	// Route set
	RouteSet []string

	// Contact addresses
	LocalContact  string
	RemoteContact string

	// Dialog state
	State SessionState

	// SDP
	LocalSDP  string
	RemoteSDP string

	// Component handling this session
	Component string // "pcscf", "icscf", "scscf"
}

// SessionState represents the state of a session
type SessionState string

const (
	SessionStateInit       SessionState = "init"
	SessionStateProceeding SessionState = "proceeding"
	SessionStateEarly      SessionState = "early"
	SessionStateConfirmed  SessionState = "confirmed"
	SessionStateTerminated SessionState = "terminated"
)

// Registration represents a subscriber registration
type Registration struct {
	IMPI        string
	IMPU        string
	Contact     string
	Expires     int
	Path        []string
	SCSCFName   string
	State       RegistrationState
}

// RegistrationState represents the state of a registration
type RegistrationState string

const (
	RegistrationStateInit       RegistrationState = "init"
	RegistrationStateRegistered RegistrationState = "registered"
	RegistrationStateUnregistered RegistrationState = "unregistered"
)
