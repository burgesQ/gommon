package webtest

import (
	"io"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/stretchr/testify/require"
)

type Case struct {
	verb     string
	path     string
	body     string
	contains string
	what     string
	payload  []byte
	headers  [][2]string
	code     int
}

// ensure type implement interface at compile time
var _ TestCase = (*Case)(nil)

func Get200() *Case {
	return (&Case{}).Code(http.StatusOK).Get().(*Case) //nolint: forcetypeassert
}

func Run(t *testing.T, uri string, tc TestCaseRun) {
	t.Helper()

	path, payload, verb := tc.GetPath(), tc.GetPayload(), tc.GetVerb()

	t.Logf("\t\t [?] running %s", tc.GetWhat())

	switch verb {
	case http.MethodGet:
		t.Logf("\t\t\t~~ GET %q", uri+path)
		webtest.RequestAndTestAPI(t, uri+path,
			func(t *testing.T, resp *http.Response) {
				t.Helper()
				Check(t, tc, resp)
			})

	case http.MethodPost:
		t.Logf("\t\t\t~~ POST %q %q", uri+path, payload)
		webtest.PushAndTestAPI(t, uri+path, payload,
			func(t *testing.T, resp *http.Response) {
				t.Helper()
				Check(t, tc, resp)
			})

	case http.MethodDelete:
		t.Logf("\t\t\t~~ DELETE %q %q", uri+path, payload)
		webtest.DeleteAndTestAPI(t, uri+path,
			func(t *testing.T, resp *http.Response) {
				t.Helper()
				Check(t, tc, resp)
			})

	default:
		t.Errorf("no verb specified for test")
	}
}

func Check(t *testing.T, tc TestCaseChecker, resp *http.Response) {
	t.Helper()

	getB := func() []byte {
		body, err := io.ReadAll(resp.Body)
		require.Nil(t, err)

		defer resp.Body.Close()
		t.Logf("\n\t\t\t [!] recv [%d] - [%s]\n\n", resp.StatusCode, body)

		return body
	}

	t.Logf("\t\t\t\t~~ testing request body\n")
	{
		switch b, bc := tc.GetBody(), tc.GetContains(); {
		case b != "":
			webtest.BodyStr(t, b, getB())
		case bc != "":
			webtest.BodyContainsStr(t, bc, getB())
		default:
			t.Log("~~ no check run against request ~~")
		}
	}

	t.Logf("\t\t\t\t~~ testing request status code\n")
	{
		webtest.StatusCode(t, tc.GetCode(), resp)
	}

	t.Logf("\t\t\t\t~~ testing request headers\n")
	{
		if h := tc.GetHeaders(); len(h) > 0 {
			webtest.Headers(t, resp, h...)
		}
	}
}

type TestCaseSetter interface {
	Body(string) TestCase
	Code(int) TestCase
	Contains(string) TestCase
	Delete() TestCase
	Get() TestCase
	Headers([][2]string) TestCase
	HeaderAdd([2]string) TestCase
	Path(string) TestCase
	PathAdd(string) TestCase
	Payload([]byte) TestCase
	PayloadStr(string) TestCase
	Post() TestCase
	What(w string) TestCase
}

func (tc *Case) Body(b string) TestCase         { tc.body = b; return tc }
func (tc *Case) Code(v int) TestCase            { tc.code = v; return tc }
func (tc *Case) Contains(c string) TestCase     { tc.contains = c; return tc }
func (tc *Case) Delete() TestCase               { tc.verb = http.MethodDelete; return tc }
func (tc *Case) Get() TestCase                  { tc.verb = http.MethodGet; return tc }
func (tc *Case) Headers(h [][2]string) TestCase { tc.headers = h; return tc }
func (tc *Case) HeaderAdd(h [2]string) TestCase { tc.headers = append(tc.headers, h); return tc }
func (tc *Case) Path(p string) TestCase         { tc.path = p; return tc }
func (tc *Case) PathAdd(p string) TestCase      { tc.path += p; return tc }
func (tc *Case) Payload(b []byte) TestCase      { tc.payload = b; return tc }
func (tc *Case) PayloadStr(b string) TestCase   { return tc.Payload([]byte(b)) }
func (tc *Case) Post() TestCase                 { tc.verb = http.MethodPost; return tc }
func (tc *Case) What(w string) TestCase         { tc.what = w; return tc }

type TestCaseChecker interface {
	GetCode() int
	GetBody() string
	GetContains() string
	GetHeaders() [][2]string
}

func (tc *Case) GetCode() int            { return tc.code }
func (tc *Case) GetBody() string         { return tc.body }
func (tc *Case) GetContains() string     { return tc.contains }
func (tc *Case) GetHeaders() [][2]string { return tc.headers }

type TestCaseRunner interface {
	GetPath() string
	GetPayload() []byte
	GetWhat() string
	GetVerb() string
}

func (tc *Case) GetPath() string    { return tc.path }
func (tc *Case) GetPayload() []byte { return tc.payload }
func (tc *Case) GetVerb() string    { return tc.verb }
func (tc *Case) GetWhat() string    { return tc.what }

type TestCaseRun interface {
	TestCaseChecker
	TestCaseRunner
}

type TestCase interface {
	TestCaseChecker
	TestCaseRunner
	TestCaseSetter
}
