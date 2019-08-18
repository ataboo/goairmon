package testhelpers

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

type FakeContext struct {
	Values     map[string]interface{}
	FakeWriter *httptest.ResponseRecorder
	response   *echo.Response
}

// Request returns `*http.Request`.
func (c *FakeContext) Request() *http.Request {
	return httptest.NewRequest("GET", "http://test.com", strings.NewReader("request body"))
}

// SetRequest sets `*http.Request`.
func (c *FakeContext) SetRequest(r *http.Request) {
	panic("not implemented")
}

// Response returns `*Response`.
func (c *FakeContext) Response() *echo.Response {
	if c.FakeWriter == nil {
		panic("not implemented")
	}

	if c.response == nil {
		c.response = &echo.Response{
			Writer: c.FakeWriter,
		}
	}

	return c.response
}

// IsTLS returns true if HTTP connection is TLS otherwise false.
func (c *FakeContext) IsTLS() bool {
	panic("not implemented")
}

// IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
func (c *FakeContext) IsWebSocket() bool {
	panic("not implemented")
}

// Scheme returns the HTTP protocol scheme, `http` or `https`.
func (c *FakeContext) Scheme() string {
	panic("not implemented")
}

// RealIP returns the client's network address based on `X-Forwarded-For`
// or `X-Real-IP` request header.
func (c *FakeContext) RealIP() string {
	panic("not implemented")
}

// Path returns the registered path for the handler.
func (c *FakeContext) Path() string {
	panic("not implemented")
}

// SetPath sets the registered path for the handler.
func (c *FakeContext) SetPath(p string) {
	panic("not implemented")
}

// Param returns path parameter by name.
func (c *FakeContext) Param(name string) string {
	panic("not implemented")
}

// ParamNames returns path parameter names.
func (c *FakeContext) ParamNames() []string {
	panic("not implemented")
}

// SetParamNames sets path parameter names.
func (c *FakeContext) SetParamNames(names ...string) {
	panic("not implemented")
}

// ParamValues returns path parameter values.
func (c *FakeContext) ParamValues() []string {
	panic("not implemented")
}

// SetParamValues sets path parameter values.
func (c *FakeContext) SetParamValues(values ...string) {
	panic("not implemented")
}

// QueryParam returns the query param for the provided name.
func (c *FakeContext) QueryParam(name string) string {
	panic("not implemented")
}

// QueryParams returns the query parameters as `url.Values`.
func (c *FakeContext) QueryParams() url.Values {
	panic("not implemented")
}

// QueryString returns the URL query string.
func (c *FakeContext) QueryString() string {
	panic("not implemented")
}

// FormValue returns the form field value for the provided name.
func (c *FakeContext) FormValue(name string) string {
	panic("not implemented")
}

// FormParams returns the form parameters as `url.Values`.
func (c *FakeContext) FormParams() (url.Values, error) {
	panic("not implemented")
}

// FormFile returns the multipart form file for the provided name.
func (c *FakeContext) FormFile(name string) (*multipart.FileHeader, error) {
	panic("not implemented")
}

// MultipartForm returns the multipart form.
func (c *FakeContext) MultipartForm() (*multipart.Form, error) {
	panic("not implemented")
}

// Cookie returns the named cookie provided in the request.
func (c *FakeContext) Cookie(name string) (*http.Cookie, error) {
	panic("not implemented")
}

// SetCookie adds a `Set-Cookie` header in HTTP response.
func (c *FakeContext) SetCookie(cookie *http.Cookie) {
	panic("not implemented")
}

// Cookies returns the HTTP cookies sent with the request.
func (c *FakeContext) Cookies() []*http.Cookie {
	panic("not implemented")
}

// Get retrieves data from the context.
func (c *FakeContext) Get(key string) interface{} {
	val, ok := c.Values[key]
	if !ok {
		return nil
	}

	return val
}

// Set saves data in the context.
func (c *FakeContext) Set(key string, val interface{}) {
	c.Values[key] = val
}

// Bind binds the request body into provided type `i`. The default binder
// does it based on Content-Type header.
func (c *FakeContext) Bind(i interface{}) error {
	panic("not implemented")
}

// Validate validates provided `i`. It is usually called after `Context#Bind()`.
// Validator must be registered using `Echo#Validator`.
func (c *FakeContext) Validate(i interface{}) error {
	panic("not implemented")
}

// Render renders a template with data and sends a text/html response with status
// code. Renderer must be registered using `Echo.Renderer`.
func (c *FakeContext) Render(code int, name string, data interface{}) error {
	panic("not implemented")
}

// HTML sends an HTTP response with status code.
func (c *FakeContext) HTML(code int, html string) error {
	panic("not implemented")
}

// HTMLBlob sends an HTTP blob response with status code.
func (c *FakeContext) HTMLBlob(code int, b []byte) error {
	panic("not implemented")
}

// String sends a string response with status code.
func (c *FakeContext) String(code int, s string) error {
	c.Response().Status = code
	c.Response().Writer.Write([]byte(s))

	return nil
}

// JSON sends a JSON response with status code.
func (c *FakeContext) JSON(code int, i interface{}) error {
	panic("not implemented")
}

// JSONPretty sends a pretty-print JSON with status code.
func (c *FakeContext) JSONPretty(code int, i interface{}, indent string) error {
	panic("not implemented")
}

// JSONBlob sends a JSON blob response with status code.
func (c *FakeContext) JSONBlob(code int, b []byte) error {
	panic("not implemented")
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (c *FakeContext) JSONP(code int, callback string, i interface{}) error {
	panic("not implemented")
}

// JSONPBlob sends a JSONP blob response with status code. It uses `callback`
// to construct the JSONP payload.
func (c *FakeContext) JSONPBlob(code int, callback string, b []byte) error {
	panic("not implemented")
}

// XML sends an XML response with status code.
func (c *FakeContext) XML(code int, i interface{}) error {
	panic("not implemented")
}

// XMLPretty sends a pretty-print XML with status code.
func (c *FakeContext) XMLPretty(code int, i interface{}, indent string) error {
	panic("not implemented")
}

// XMLBlob sends an XML blob response with status code.
func (c *FakeContext) XMLBlob(code int, b []byte) error {
	panic("not implemented")
}

// Blob sends a blob response with status code and content type.
func (c *FakeContext) Blob(code int, contentType string, b []byte) error {
	panic("not implemented")
}

// Stream sends a streaming response with status code and content type.
func (c *FakeContext) Stream(code int, contentType string, r io.Reader) error {
	panic("not implemented")
}

// File sends a response with the content of the file.
func (c *FakeContext) File(file string) error {
	panic("not implemented")
}

// Attachment sends a response as attachment, prompting client to save the
// file.
func (c *FakeContext) Attachment(file string, name string) error {
	panic("not implemented")
}

// Inline sends a response as inline, opening the file in the browser.
func (c *FakeContext) Inline(file string, name string) error {
	panic("not implemented")
}

// NoContent sends a response with no body and a status code.
func (c *FakeContext) NoContent(code int) error {
	panic("not implemented")
}

// Redirect redirects the request to a provided URL with status code.
func (c *FakeContext) Redirect(code int, url string) error {
	panic("not implemented")
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *FakeContext) Error(err error) {
	panic("not implemented")
}

// Handler returns the matched handler by router.
func (c *FakeContext) Handler() echo.HandlerFunc {
	panic("not implemented")
}

// SetHandler sets the matched handler by router.
func (c *FakeContext) SetHandler(h echo.HandlerFunc) {
	panic("not implemented")
}

// Logger returns the `Logger` instance.
func (c *FakeContext) Logger() echo.Logger {
	panic("not implemented")
}

// Echo returns the `Echo` instance.
func (c *FakeContext) Echo() *echo.Echo {
	panic("not implemented")
}

// Reset resets the context after request completes. It must be called along
// with `Echo#AcquireContext()` and `Echo#ReleaseContext()`.
// See `Echo#ServeHTTP()`
func (c *FakeContext) Reset(r *http.Request, w http.ResponseWriter) {
	panic("not implemented")
}

type FakeCookieStore struct {
	Sessions map[string]*sessions.Session
}

// Get should return a cached session.
func (s *FakeCookieStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	sess, ok := s.Sessions[name]
	if !ok {
		sess, _ = s.New(r, name)
	}

	return sess, nil
}

// New should create and return a new session.
//
// Note that New should never return a nil session, even in the case of
// an error if using the Registry infrastructure to cache the session.
func (s *FakeCookieStore) New(r *http.Request, name string) (*sessions.Session, error) {
	sess := sessions.NewSession(s, name)
	s.Sessions[name] = sess

	return sess, nil
}

// Save should persist session to the underlying store implementation.
func (s *FakeCookieStore) Save(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	sess.IsNew = false
	s.Sessions[sess.Name()] = sess

	return nil
}

func EmptyHandler(ctx echo.Context) error {
	return nil
}
