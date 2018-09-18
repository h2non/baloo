package baloo

import (
	"io"
	"net/http"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
	"gopkg.in/h2non/gentleman.v2/plugins/multipart"
)

var (
	// UserAgent represents the static user agent name and version.
	UserAgent = "baloo/" + Version
)

// Request HTTP entity for gentleman.
// Provides middleware capabilities, built-in context
// and convenient methods to easily setup request params.
type Request struct {
	// tested stores if the request was already tested.
	tested bool

	// Optional reference to the parent Client instance.
	Client *Client

	// Request stores the reference to gentleman.Request instance.
	Request *gentleman.Request
}

// NewRequest creates a new Request entity.
func NewRequest() *Request {
	req := gentleman.NewRequest()
	req.SetHeader("User-Agent", UserAgent)
	return &Request{Request: req}
}

// SetClient Attach a client to the current Request
// This is mostly done internally.
func (r *Request) SetClient(cli *Client) *Request {
	r.Client = cli
	r.Request.SetClient(cli.Client)
	return r
}

// Method defines the HTTP verb to be used.
func (r *Request) Method(method string) *Request {
	r.Request.Middleware.UseRequest(func(ctx *context.Context, h context.Handler) {
		ctx.Request.Method = method
		h.Next(ctx)
	})
	return r
}

// URL parses and defines the URL to be used in the outgoing request.
func (r *Request) URL(uri string) *Request {
	r.Request.URL(uri)
	return r
}

// BaseURL parses the given URL and uses the URL schema and host in the outgoing request.
func (r *Request) BaseURL(uri string) *Request {
	r.Request.BaseURL(uri)
	return r
}

// Path defines the request URL path to be used in the outgoing request.
func (r *Request) Path(path string) *Request {
	r.Request.Path(path)
	return r
}

// AddPath defines the request URL path to be used in the outgoing request.
func (r *Request) AddPath(path string) *Request {
	r.Request.AddPath(path)
	return r
}

// Param replaces a path param based on the given param name and value.
func (r *Request) Param(name, value string) *Request {
	r.Request.Param(name, value)
	return r
}

// Params replaces path params based on the given params key-value map.
func (r *Request) Params(params map[string]string) *Request {
	r.Request.Params(params)
	return r
}

// SetQuery sets a new URL query param field.
// If another query param exists with the same key, it will be overwritten.
func (r *Request) SetQuery(name, value string) *Request {
	r.Request.SetQuery(name, value)
	return r
}

// AddQuery adds a new URL query param field
// without overwriting any existent query field.
func (r *Request) AddQuery(name, value string) *Request {
	r.Request.AddQuery(name, value)
	return r
}

// SetQueryParams sets URL query params based on the given map.
func (r *Request) SetQueryParams(params map[string]string) *Request {
	r.Request.SetQueryParams(params)
	return r
}

// SetHeader sets a new header field by name and value.
// If another header exists with the same key, it will be overwritten.
func (r *Request) SetHeader(name, value string) *Request {
	r.Request.SetHeader(name, value)
	return r
}

// AddHeader adds a new header field by name and value
// without overwriting any existent header.
func (r *Request) AddHeader(name, value string) *Request {
	r.Request.AddHeader(name, value)
	return r
}

// SetHeaders adds new header fields based on the given map.
func (r *Request) SetHeaders(fields map[string]string) *Request {
	r.Request.SetHeaders(fields)
	return r
}

// AddCookie sets a new cookie field bsaed on the given http.Cookie struct
// without overwriting any existent cookie.
func (r *Request) AddCookie(cookie *http.Cookie) *Request {
	r.Request.AddCookie(cookie)
	return r
}

// AddCookies sets a new cookie field based on a list of http.Cookie
// without overwriting any existent cookie.
func (r *Request) AddCookies(data []*http.Cookie) *Request {
	r.Request.AddCookies(data)
	return r
}

// CookieJar creates a cookie jar to store HTTP cookies when they are sent down.
func (r *Request) CookieJar() *Request {
	r.Request.CookieJar()
	return r
}

// Type defines the Content-Type header field based on the given type name alias or value.
// You can use the following content type aliases: json, xml, form, html, text and urlencoded.
func (r *Request) Type(name string) *Request {
	r.Request.Type(name)
	return r
}

// Body defines the request body based on a io.Reader stream.
func (r *Request) Body(reader io.Reader) *Request {
	r.Request.Body(reader)
	return r
}

// BodyString defines the request body based on the given string.
// If using this method, you should define the proper Content-Type header
// representing the real content MIME type.
func (r *Request) BodyString(data string) *Request {
	r.Request.BodyString(data)
	return r
}

// JSON serializes and defines as request body based on the given input.
// The proper Content-Type header will be transparently added for you.
func (r *Request) JSON(data interface{}) *Request {
	r.Request.JSON(data)
	return r
}

// XML serializes and defines the request body based on the given input.
// The proper Content-Type header will be transparently added for you.
func (r *Request) XML(data interface{}) *Request {
	r.Request.XML(data)
	return r
}

// Form serializes and defines the request body as multipart/form-data
// based on the given form data.
func (r *Request) Form(data multipart.FormData) *Request {
	r.Request.Form(data)
	return r
}

// File serializes and defines the request body as multipart/form-data
// containing one file field.
func (r *Request) File(name string, reader io.Reader) *Request {
	r.Request.File(name, reader)
	return r
}

// Files serializes and defines the request body as multipart/form-data
// containing the given file fields.
func (r *Request) Files(files []multipart.FormFile) *Request {
	r.Request.Files(files)
	return r
}

// Send executes the current request and returns
// the response or error.
func (r *Request) Send() (*gentleman.Response, error) {
	return r.Request.Send()
}

// Expect creates and returns the request test expectation suite.
func (r *Request) Expect(t TestingT) *Expect {
	if r.tested {
		t.Error("[baloo] request already tested")
		return nil
	}
	r.tested = true
	return NewExpect(r).BindTest(t)
}

// Assert is an alias to .Expect().
func (r *Request) Assert(t TestingT) *Expect {
	return r.Expect(t)
}

// Use uses a new plugin in the middleware stack.
func (r *Request) Use(p plugin.Plugin) *Request {
	r.Request.Use(p)
	return r
}

// UseRequest uses a request middleware handler.
func (r *Request) UseRequest(fn context.HandlerFunc) *Request {
	r.Request.UseRequest(fn)
	return r
}

// UseResponse uses a response middleware handler.
func (r *Request) UseResponse(fn context.HandlerFunc) *Request {
	r.Request.UseResponse(fn)
	return r
}

// UseError uses an error middleware handler.
func (r *Request) UseError(fn context.HandlerFunc) *Request {
	r.Request.UseError(fn)
	return r
}

// UseHandler uses an new middleware handler for the given phase.
func (r *Request) UseHandler(phase string, fn context.HandlerFunc) *Request {
	r.Request.UseHandler(phase, fn)
	return r
}

// Clone creates a new side-effects free Request based on the current one.
func (r *Request) Clone() *Request {
	return &Request{Request: r.Request.Clone()}
}
