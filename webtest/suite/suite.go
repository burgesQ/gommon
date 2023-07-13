package webtest

import (
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
)

type Case struct {
	verb     string
	path     string
	body     string
	contains string
	payload  []byte
	code     int
}

// ensure type implement interface at compile time
var _ TestCase = (*Case)(nil)

func Get200() *Case {
	return (&Case{}).Code(http.StatusOK).Get().(*Case) //nolint: forcetypeassert
}

func Run(t *testing.T, uri string, tc TestCase) {
	t.Helper()

	path, payload, verb := tc.GetPath(), tc.GetPayload(), tc.GetVerb()

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

	default:
		t.Errorf("no verb specified for test")
	}
}

func Check(t *testing.T, tc TestCaseChecker, resp *http.Response) {
	t.Helper()

	t.Log("\t\t\t\t~~ testing request status code")
	{
		webtest.StatusCode(t, tc.GetCode(), resp)
	}

	t.Log("\t\t\t\t~~ testing request body code")
	{
		switch b, bc := tc.GetBody(), tc.GetContains(); {
		case b != "":
			webtest.Body(t, b, resp)
		case bc != "":
			webtest.BodyContains(t, bc, resp)
		default:
			t.Log("~~ no check run against request body ~~")
		}
	}
}

type TestCaseSetter interface {
	Body(string) TestCase
	Code(int) TestCase
	Contains(string) TestCase
	Path(string) TestCase
	PathAdd(string) TestCase
	Payload([]byte) TestCase
	Post() TestCase
	Get() TestCase
}

func (tc *Case) Body(b string) TestCase {
	tc.body = b

	return tc
}

func (tc *Case) Code(v int) TestCase {
	tc.code = v

	return tc
}

func (tc *Case) Contains(c string) TestCase {
	tc.contains = c

	return tc
}

func (tc *Case) Path(p string) TestCase {
	tc.path = p

	return tc
}

func (tc *Case) PathAdd(p string) TestCase {
	tc.path += p

	return tc
}

func (tc *Case) Payload(b []byte) TestCase {
	tc.payload = b

	return tc
}

func (tc *Case) Post() TestCase {
	tc.verb = http.MethodPost

	return tc
}

func (tc *Case) Get() TestCase {
	tc.verb = http.MethodGet

	return tc
}

type TestCaseChecker interface {
	GetCode() int
	GetBody() string
	GetContains() string
}

func (tc *Case) GetCode() int        { return tc.code }
func (tc *Case) GetBody() string     { return tc.body }
func (tc *Case) GetContains() string { return tc.contains }

type TestCaseRunner interface {
	GetPath() string
	GetPayload() []byte
	GetVerb() string
}

func (tc *Case) GetPath() string    { return tc.path }
func (tc *Case) GetPayload() []byte { return tc.payload }
func (tc *Case) GetVerb() string    { return tc.verb }

type TestCase interface {
	TestCaseChecker
	TestCaseRunner
	TestCaseSetter
}
