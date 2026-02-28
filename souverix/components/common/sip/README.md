# SIP Framework

Common SIP (Session Initiation Protocol) framework for Souverix components.

## Features

- **SIP Message Parsing**: Parse and construct SIP messages
- **Handler Framework**: Register handlers for SIP methods (INVITE, REGISTER, etc.)
- **Middleware Support**: Add middleware for logging, validation, header manipulation
- **Component Framework**: High-level component abstraction
- **Proxy Support**: Built-in proxy/forward functionality

## Usage

### Basic Component Setup

```go
import "github.com/dasmlab/souverix/common/sip"

// Create component
component := sip.NewComponent("pcscf", "localhost", 8081, logger)

// Setup default middleware
component.SetupDefaultMiddleware()

// Register SIP method handlers
component.RegisterMethod("INVITE", func(ctx context.Context, msg *sip.Message, c *gin.Context) (*sip.Message, error) {
    // Handle INVITE
    response := sip.CreateResponse(200, "OK", msg)
    return response, nil
})

component.RegisterMethod("REGISTER", func(ctx context.Context, msg *sip.Message, c *gin.Context) (*sip.Message, error) {
    // Handle REGISTER
    response := sip.CreateResponse(200, "OK", msg)
    return response, nil
})

// Register routes
component.RegisterRoutes("/sip")

// Start server
component.StartAsync()
```

### Parsing SIP Messages

```go
rawSIP := `INVITE sip:user@example.com SIP/2.0
Via: SIP/2.0/UDP host.example.com;branch=z9hG4bK776asdhds
From: <sip:caller@example.com>;tag=1928301774
To: <sip:user@example.com>
Call-ID: a84b4c76e66710@host.example.com
CSeq: 1 INVITE
Content-Length: 0
`

msg, err := sip.ParseMessage(rawSIP)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Method: %s\n", msg.Method)
fmt.Printf("URI: %s\n", msg.URI)
fmt.Printf("Call-ID: %s\n", msg.GetCallID())
```

### Creating SIP Messages

```go
// Create request
request := sip.CreateRequest("INVITE", "sip:user@example.com")
request.SetHeader("From", "<sip:caller@example.com>")
request.SetHeader("To", "<sip:user@example.com>")
request.SetHeader("Call-ID", "unique-call-id")
request.SetHeader("CSeq", "1 INVITE")

// Create response
response := sip.CreateResponse(200, "OK", request)
response.SetHeader("Contact", "<sip:server@example.com>")
```

### Forwarding/Proxying

```go
// Forward to another component
response, err := component.Forward(ctx, msg, "icscf.example.com", 8082)

// Proxy request (adds Record-Route)
response, err := component.ProxyRequest(ctx, msg, "scscf.example.com", 8083)
```

### Middleware

```go
// Add custom middleware
component.Use(func(ctx context.Context, msg *sip.Message, c *gin.Context) (context.Context, error) {
    // Custom processing
    return ctx, nil
})
```

## Integration with Components

Components should integrate SIP handling into their main server (r1):

```go
// In main.go
sipComponent := sip.NewComponent("pcscf", "localhost", 8081, logger)
sipComponent.SetupDefaultMiddleware()

// Register handlers
sipComponent.RegisterMethod("INVITE", handleINVITE)
sipComponent.RegisterMethod("REGISTER", handleREGISTER)

// Register on main router (r1)
sipComponent.RegisterRoutes("/sip")

// Start SIP component
sipComponent.StartAsync()
```

## Supported SIP Methods

- INVITE
- REGISTER
- ACK
- BYE
- CANCEL
- OPTIONS
- UPDATE
- PRACK
- INFO
- REFER
- NOTIFY
- SUBSCRIBE

## Headers

Common SIP headers are supported:
- Via
- From
- To
- Call-ID
- CSeq
- Contact
- Record-Route
- Route
- Max-Forwards
- Content-Type
- Content-Length
