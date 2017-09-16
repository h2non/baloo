package baloo

import (
	"net/http"

	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

// NewHandler is a convenient alias to gentleman.NewHandler factory.
var NewHandler = gentleman.NewHandler

// Client represents a high-level HTTP client entity capable
// with a built-in middleware and context.
type Client struct {
	// Parent stores an optional parent baloo Client instance.
	Parent *Client
	// Client entity has it's own Context that will be inherited by requests or child clients.
	Client *gentleman.Client
}

// New creates a new high level client entity
// able to perform HTTP requests.
func New(url string) *Client {
	cli := gentleman.New()
	cli.URL(url)
	return &Client{Client: cli}
}

// Request creates a new Request based on the current Client
func (c *Client) Request() *Request {
	req := NewRequest()
	req.SetClient(c)
	return req
}

// Get creates a new GET request.
func (c *Client) Get(path string) *Request {
	req := c.Request()
	req.Method("GET")
	req.Path(path)
	return req
}

// Post creates a new POST request.
func (c *Client) Post(path string) *Request {
	req := c.Request()
	req.Method("POST")
	req.Path(path)
	return req
}

// Put creates a new PUT request.
func (c *Client) Put(path string) *Request {
	req := c.Request()
	req.Method("PUT")
	req.Path(path)
	return req
}

// Delete creates a new DELETE request.
func (c *Client) Delete(path string) *Request {
	req := c.Request()
	req.Method("DELETE")
	req.Path(path)
	return req
}

// Patch creates a new PATCH request.
func (c *Client) Patch(path string) *Request {
	req := c.Request()
	req.Method("PATCH")
	req.Path(path)
	return req
}

// Head creates a new HEAD request.
func (c *Client) Head(path string) *Request {
	req := c.Request()
	req.Method("HEAD")
	req.Path(path)
	return req
}

// Method defines a the default HTTP method used by outgoing client requests.
func (c *Client) Method(name string) *Client {
	c.Client.Middleware.UseRequest(func(ctx *context.Context, h context.Handler) {
		ctx.Request.Method = name
		h.Next(ctx)
	})
	return c
}

// URL defines the URL for client requests.
// Useful to define at client level the base URL and base path used by child requests.
func (c *Client) URL(uri string) *Client {
	c.Client.URL(uri)
	return c
}

// BaseURL defines the URL schema and host for client requests.
// Useful to define at client level the base URL used by client child requests.
func (c *Client) BaseURL(uri string) *Client {
	c.Client.BaseURL(uri)
	return c
}

// Path defines the URL base path for client requests.
func (c *Client) Path(path string) *Client {
	c.Client.Path(path)
	return c
}

// Param replaces a path param based on the given param name and value.
func (c *Client) Param(name, value string) *Client {
	c.Client.Param(name, value)
	return c
}

// Params replaces path params based on the given params key-value map.
func (c *Client) Params(params map[string]string) *Client {
	c.Client.Params(params)
	return c
}

// SetHeader sets a new header field by name and value.
// If another header exists with the same key, it will be overwritten.
func (c *Client) SetHeader(key, value string) *Client {
	c.Client.SetHeader(key, value)
	return c
}

// AddHeader adds a new header field by name and value
// without overwriting any existent header.
func (c *Client) AddHeader(name, value string) *Client {
	c.Client.AddHeader(name, value)
	return c
}

// SetHeaders adds new header fields based on the given map.
func (c *Client) SetHeaders(fields map[string]string) *Client {
	c.Client.SetHeaders(fields)
	return c
}

// AddCookie sets a new cookie field bsaed on the given http.Cookie struct
// without overwriting any existent cookie.
func (c *Client) AddCookie(cookie *http.Cookie) *Client {
	c.Client.AddCookie(cookie)
	return c
}

// AddCookies sets a new cookie field based on a list of http.Cookie
// without overwriting any existent cookie.
func (c *Client) AddCookies(data []*http.Cookie) *Client {
	c.Client.AddCookies(data)
	return c
}

// CookieJar creates a cookie jar to store HTTP cookies when they are sent down.
func (c *Client) CookieJar() *Client {
	c.Client.CookieJar()
	return c
}

// Use uses a new plugin to the middleware stack.
func (c *Client) Use(p plugin.Plugin) *Client {
	c.Client.Use(p)
	return c
}

// UseRequest uses a new middleware function for request phase.
func (c *Client) UseRequest(fn context.HandlerFunc) *Client {
	c.Client.UseRequest(fn)
	return c
}

// UseResponse uses a new middleware function for response phase.
func (c *Client) UseResponse(fn context.HandlerFunc) *Client {
	c.Client.UseResponse(fn)
	return c
}

// UseError uses a new middleware function for error phase.
func (c *Client) UseError(fn context.HandlerFunc) *Client {
	c.Client.UseError(fn)
	return c
}

// UseHandler uses a new middleware function for the given phase.
func (c *Client) UseHandler(phase string, fn context.HandlerFunc) *Client {
	c.Client.UseHandler(phase, fn)
	return c
}

// UseParent uses another Client as parent
// inheriting its middleware stack and configuration.
func (c *Client) UseParent(parent *Client) *Client {
	c.Parent = parent
	c.Client.UseParent(parent.Client)
	return c
}
