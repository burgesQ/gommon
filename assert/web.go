package assert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type (
	// HandlerForTest implement the function signature used to check the req/resp
	HandlerForTest = func(t *testing.T, resp *http.Response)
)

// Body fetch and assert that the body of the http.Response is the same than expected
func Body(t *testing.T, resp *http.Response, expected string) {
	t.Helper()
	StringEqual(t, fetchBody(t, resp), expected)
}

// BodyDiffere fetch and assert that the body of the http.Response differ than expected
func BodyDiffere(t *testing.T, resp *http.Response, expected string) {
	t.Helper()
	StringNotEqual(t, fetchBody(t, resp), expected)
}

// StatusCode assert the status code of the response
func StatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	IntEqual(t, resp.StatusCode, expected)
}

// Header assert value of the given header key:vak in the htt.Response param
func Header(t *testing.T, resp *http.Response, key, val string) bool {
	// test existence
	if out, ok := resp.Header[key]; !ok || len(out) == 0 || out[0] != val {
		t.Errorf("Invalid response header [%s] expected: [%s]", out[0], val)
		return false
	}

	return true
}

// RequestAndTestAPI request an API then run the test handler
func RequestAndTestAPI(t *testing.T, url string, handler HandlerForTest) {
	t.Helper()
	var resp = requestAPI(t, url)
	defer resp.Body.Close()
	handler(t, resp)
}

// PushAndTestAPI post to an API then run the test handler
// The sub method try to send an `application/json` encoded content
func PushAndTestAPI(t *testing.T, path string, content []byte, handler HandlerForTest) {
	t.Helper()
	var resp = pushAPI(t, path, content)
	defer resp.Body.Close()
	handler(t, resp)
}

func fetchBody(t *testing.T, resp *http.Response) string {
	var tmp, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error fetching the body response : %s", err.Error())
	}
	defer resp.Body.Close()
	return string(tmp)
}

func requestAPI(t *testing.T, url string) *http.Response {
	var resp, err = http.Get(url)
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	return resp
}

func pushAPI(t *testing.T, url string, content []byte) *http.Response {
	var req, err = http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
		t.Fatalf("can't post the new request : %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	return resp
}
