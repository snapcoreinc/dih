package handler

import (
	"context"
	"net/http"
)

// Response of module call
type Response struct {

	// Body the body will be written back
	Body []byte

	// StatusCode needs to be populated with value such as http.StatusOK
	StatusCode int

	// Header is optional and contains any additional headers the module response should set
	Header http.Header
}

// Request of module call
type Request struct {
	Body        []byte
	Header      http.Header
	QueryString string
	Method      string
	Host        string
	ctx         context.Context
}

type Context struct {
	moduleName         string
	moduleVersion      string
	identity           string
	roles              string
	isReqBase64Encoded bool
}

// Context is set for optional cancellation inflight requests.
func (r *Request) Context() context.Context {
	return r.ctx
}

// WithContext overides the context for the Request struct
func (r *Request) WithContext(ctx context.Context) {
	// AE: Not keen on panic mid-flow in user-code, however stdlib also appears to do
	// this. https://golang.org/src/net/http/request.go
	// This is not setting a precedent for broader use of "panic" to handle errors.
	if ctx == nil {
		panic("nil context")
	}
	r.ctx = ctx
}

// ModuleHandler used for a serverless Go method invocation
type ModuleHandler interface {
	handleRequest(ctx Context, req Request) (Response, error)
}
