package middleware

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"net/http/httptest"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

func TestServiceProvider(t *testing.T) {
	provider := NewServiceProvider()

	provider.Register("first_key", "first_value")
	provider.Register("second_key", "second_value")

	ctx := &_fakeContext{values: make(map[string]interface{})}

	provider.BindMiddleware(ctx)(nil)

	if ctx.Get("first_key").(string) != "first_value" {
		t.Error("expected bound value")
	}

	if ctx.Get("second_key").(string) != "second_value" {
		t.Error("expected bound value")
	}

}

type _fakeContext struct {
	values     map[string]interface{}
	fakeWriter *httptest.ResponseRecorder
}

// Request returns `*http.Request`.
func (c *_fakeContext) Request() *http.Request {
	return httptest.NewRequest("GET", "http://test.com", strings.NewReader("request body"))
}

// SetRequest sets `*http.Request`.
func (c *_fakeContext) SetRequest(r *http.Request) {
	panic("not implemented")
}

// Response returns `*Response`.
func (c *_fakeContext) Response() *echo.Response {
	if c.fakeWriter == nil {
		panic("not implemented")
	}

	return &echo.Response{
		Writer: c.fakeWriter,
	}
}

// IsTLS returns true if HTTP connection is TLS otherwise false.
func (c *_fakeContext) IsTLS() bool {
	panic("not implemented")
}

// IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
func (c *_fakeContext) IsWebSocket() bool {
	panic("not implemented")
}

// Scheme returns the HTTP protocol scheme, `http` or `https`.
func (c *_fakeContext) Scheme() string {
	panic("not implemented")
}

// RealIP returns the client's network address based on `X-Forwarded-For`
// or `X-Real-IP` request header.
func (c *_fakeContext) RealIP() string {
	panic("not implemented")
}

// Path returns the registered path for the handler.
func (c *_fakeContext) Path() string {
	panic("not implemented")
}

// SetPath sets the registered path for the handler.
func (c *_fakeContext) SetPath(p string) {
	panic("not implemented")
}

// Param returns path parameter by name.
func (c *_fakeContext) Param(name string) string {
	panic("not implemented")
}

// ParamNames returns path parameter names.
func (c *_fakeContext) ParamNames() []string {
	panic("not implemented")
}

// SetParamNames sets path parameter names.
func (c *_fakeContext) SetParamNames(names ...string) {
	panic("not implemented")
}

// ParamValues returns path parameter values.
func (c *_fakeContext) ParamValues() []string {
	panic("not implemented")
}

// SetParamValues sets path parameter values.
func (c *_fakeContext) SetParamValues(values ...string) {
	panic("not implemented")
}

// QueryParam returns the query param for the provided name.
func (c *_fakeContext) QueryParam(name string) string {
	panic("not implemented")
}

// QueryParams returns the query parameters as `url.Values`.
func (c *_fakeContext) QueryParams() url.Values {
	panic("not implemented")
}

// QueryString returns the URL query string.
func (c *_fakeContext) QueryString() string {
	panic("not implemented")
}

// FormValue returns the form field value for the provided name.
func (c *_fakeContext) FormValue(name string) string {
	panic("not implemented")
}

// FormParams returns the form parameters as `url.Values`.
func (c *_fakeContext) FormParams() (url.Values, error) {
	panic("not implemented")
}

// FormFile returns the multipart form file for the provided name.
func (c *_fakeContext) FormFile(name string) (*multipart.FileHeader, error) {
	panic("not implemented")
}

// MultipartForm returns the multipart form.
func (c *_fakeContext) MultipartForm() (*multipart.Form, error) {
	panic("not implemented")
}

// Cookie returns the named cookie provided in the request.
func (c *_fakeContext) Cookie(name string) (*http.Cookie, error) {
	panic("not implemented")
}

// SetCookie adds a `Set-Cookie` header in HTTP response.
func (c *_fakeContext) SetCookie(cookie *http.Cookie) {
	panic("not implemented")
}

// Cookies returns the HTTP cookies sent with the request.
func (c *_fakeContext) Cookies() []*http.Cookie {
	panic("not implemented")
}

// Get retrieves data from the context.
func (c *_fakeContext) Get(key string) interface{} {
	val, ok := c.values[key]
	if !ok {
		return nil
	}

	return val
}

// Set saves data in the context.
func (c *_fakeContext) Set(key string, val interface{}) {
	c.values[key] = val
}

// Bind binds the request body into provided type `i`. The default binder
// does it based on Content-Type header.
func (c *_fakeContext) Bind(i interface{}) error {
	panic("not implemented")
}

// Validate validates provided `i`. It is usually called after `Context#Bind()`.
// Validator must be registered using `Echo#Validator`.
func (c *_fakeContext) Validate(i interface{}) error {
	panic("not implemented")
}

// Render renders a template with data and sends a text/html response with status
// code. Renderer must be registered using `Echo.Renderer`.
func (c *_fakeContext) Render(code int, name string, data interface{}) error {
	panic("not implemented")
}

// HTML sends an HTTP response with status code.
func (c *_fakeContext) HTML(code int, html string) error {
	panic("not implemented")
}

// HTMLBlob sends an HTTP blob response with status code.
func (c *_fakeContext) HTMLBlob(code int, b []byte) error {
	panic("not implemented")
}

// String sends a string response with status code.
func (c *_fakeContext) String(code int, s string) error {
	panic("not implemented")
}

// JSON sends a JSON response with status code.
func (c *_fakeContext) JSON(code int, i interface{}) error {
	panic("not implemented")
}

// JSONPretty sends a pretty-print JSON with status code.
func (c *_fakeContext) JSONPretty(code int, i interface{}, indent string) error {
	panic("not implemented")
}

// JSONBlob sends a JSON blob response with status code.
func (c *_fakeContext) JSONBlob(code int, b []byte) error {
	panic("not implemented")
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (c *_fakeContext) JSONP(code int, callback string, i interface{}) error {
	panic("not implemented")
}

// JSONPBlob sends a JSONP blob response with status code. It uses `callback`
// to construct the JSONP payload.
func (c *_fakeContext) JSONPBlob(code int, callback string, b []byte) error {
	panic("not implemented")
}

// XML sends an XML response with status code.
func (c *_fakeContext) XML(code int, i interface{}) error {
	panic("not implemented")
}

// XMLPretty sends a pretty-print XML with status code.
func (c *_fakeContext) XMLPretty(code int, i interface{}, indent string) error {
	panic("not implemented")
}

// XMLBlob sends an XML blob response with status code.
func (c *_fakeContext) XMLBlob(code int, b []byte) error {
	panic("not implemented")
}

// Blob sends a blob response with status code and content type.
func (c *_fakeContext) Blob(code int, contentType string, b []byte) error {
	panic("not implemented")
}

// Stream sends a streaming response with status code and content type.
func (c *_fakeContext) Stream(code int, contentType string, r io.Reader) error {
	panic("not implemented")
}

// File sends a response with the content of the file.
func (c *_fakeContext) File(file string) error {
	panic("not implemented")
}

// Attachment sends a response as attachment, prompting client to save the
// file.
func (c *_fakeContext) Attachment(file string, name string) error {
	panic("not implemented")
}

// Inline sends a response as inline, opening the file in the browser.
func (c *_fakeContext) Inline(file string, name string) error {
	panic("not implemented")
}

// NoContent sends a response with no body and a status code.
func (c *_fakeContext) NoContent(code int) error {
	panic("not implemented")
}

// Redirect redirects the request to a provided URL with status code.
func (c *_fakeContext) Redirect(code int, url string) error {
	panic("not implemented")
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *_fakeContext) Error(err error) {
	panic("not implemented")
}

// Handler returns the matched handler by router.
func (c *_fakeContext) Handler() echo.HandlerFunc {
	panic("not implemented")
}

// SetHandler sets the matched handler by router.
func (c *_fakeContext) SetHandler(h echo.HandlerFunc) {
	panic("not implemented")
}

// Logger returns the `Logger` instance.
func (c *_fakeContext) Logger() echo.Logger {
	panic("not implemented")
}

// Echo returns the `Echo` instance.
func (c *_fakeContext) Echo() *echo.Echo {
	panic("not implemented")
}

// Reset resets the context after request completes. It must be called along
// with `Echo#AcquireContext()` and `Echo#ReleaseContext()`.
// See `Echo#ServeHTTP()`
func (c *_fakeContext) Reset(r *http.Request, w http.ResponseWriter) {
	panic("not implemented")
}

type _fakeCookieStore struct {
	sessions map[string]*sessions.Session
}

// Get should return a cached session.
func (s *_fakeCookieStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	sess, ok := s.sessions[name]
	if !ok {
		sess, _ = s.New(r, name)
	}

	return sess, nil
}

// New should create and return a new session.
//
// Note that New should never return a nil session, even in the case of
// an error if using the Registry infrastructure to cache the session.
func (s *_fakeCookieStore) New(r *http.Request, name string) (*sessions.Session, error) {
	sess := sessions.NewSession(s, name)
	s.sessions[name] = sess

	return sess, nil
}

// Save should persist session to the underlying store implementation.
func (s *_fakeCookieStore) Save(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	sess.IsNew = false
	s.sessions[sess.Name()] = sess

	return nil
}
