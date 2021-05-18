// Package webtest import some useful method to perfome assertion for testing package.
package webtest

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// HandlerForTest implement the function signature used to check the req/resp
type HandlerForTest = func(t *testing.T, resp *http.Response)

const (
	_notEqualHeader = `assertion failed for the %q header: %q != %q`
	_ttlTest        = 5
)

// Body fetch and assert that the body of the http.Response is the same than expected
func Body(t *testing.T, resp *http.Response, expected string) {
	t.Helper()
	require.Equal(t, expected, fetchBody(t, resp))
}

// BodyDiffere fetch and assert that the body of the http.Response differ than expected
func BodyDiffere(t *testing.T, resp *http.Response, expected string) {
	t.Helper()
	require.NotEqual(t, expected, fetchBody(t, resp))
}

// StatusCode assert the status code of the response
func StatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	require.Equal(t, expected, resp.StatusCode)
}

// Header assert value of the given header key:vak in the htt.Response param
func Header(t *testing.T, resp *http.Response, key, val string) bool {
	t.Helper()
	// test existence
	if out, ok := resp.Header[key]; !ok || len(out) == 0 || out[0] != val {
		require.Equalf(t, val, out, _notEqualHeader, key, out[0], val)
		return false
	}

	return true
}

// Header assert value of the given header key:vak in the htt.Response param
func Headers(t *testing.T, resp *http.Response, kv ...[2]string) bool {
	t.Helper()

	for i := range kv {
		k, v := kv[i][0], kv[i][1]
		if !Header(t, resp, k, v) {
			return false
		}
	}

	return true
}

func DeleteAndTestAPI(t *testing.T, url string, handler HandlerForTest) {
	t.Helper()

	var resp = deleteAPI(t, url)
	defer resp.Body.Close()

	handler(t, resp)
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

func prepReq(t *testing.T, url, method string) *http.Response {
	var (
		client   = &http.Client{}
		ctx, cl  = context.WithTimeout(context.Background(), time.Second*_ttlTest)
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	)

	defer cl()

	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	return resp
}

func requestAPI(t *testing.T, url string) *http.Response {
	return prepReq(t, url, "GET")
}

func deleteAPI(t *testing.T, url string) *http.Response {
	return prepReq(t, url, "DELETE")
}

func pushAPI(t *testing.T, url string, content []byte) *http.Response {
	var (
		client   = &http.Client{}
		ctx, cl  = context.WithTimeout(context.Background(), time.Second*_ttlTest)
		req, err = http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(content))
	)

	defer cl()

	if err != nil {
		t.Fatalf("can't post the new request : %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	return resp
}
