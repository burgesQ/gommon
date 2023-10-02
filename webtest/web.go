// Package webtest import some useful method to perfome assertion for testing package.
package webtest

import (
	"bytes"
	"context"
	"io"
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
func Body(t *testing.T, expected string, resp *http.Response, msg ...any) {
	t.Helper()

	if len(msg) == 0 {
		msg = []any{"expected response body differe"}
	}

	require.Equal(t, expected, fetchBody(t, resp), msg...)
}

// Body fetch and assert that the body of the http.Response is the same than expected
func BodyStr(t *testing.T, expected string, body []byte, msg ...any) {
	t.Helper()

	if len(msg) == 0 {
		msg = []any{"expected response body differe"}
	}

	require.Equal(t, expected, string(body), msg...)
}

func BodyContains(t *testing.T, expected string, resp *http.Response, msg ...any) {
	t.Helper()

	if len(msg) == 0 {
		msg = []any{"expected response body doens't contain the substring"}
	}

	require.Contains(t, fetchBody(t, resp), expected, msg...)
}

func BodyContainsStr(t *testing.T, expected string, body []byte, msg ...any) {
	t.Helper()

	if len(msg) == 0 {
		msg = []any{"expected response body doens't contain the substring"}
	}

	require.Contains(t, string(body), expected, msg...)
}

// BodyDiffere fetch and assert that the body of the http.Response differ than expected
func BodyDiffere(t *testing.T, expected string, resp *http.Response) {
	t.Helper()
	require.NotEqual(t, expected, fetchBody(t, resp))
}

// StatusCode assert the status code of the response.
func StatusCode(t *testing.T, expected int, resp *http.Response, msg ...any) {
	t.Helper()

	if len(msg) == 0 {
		msg = []any{"expected response status code differe"}
	}

	require.Equal(t, expected, resp.StatusCode, msg...)
}

// Header assert value of the given header key:vak in the htt.Response param.
func Header(t *testing.T, key, val string, resp *http.Response) bool {
	t.Helper()
	// test existence
	if out, ok := resp.Header[key]; !ok || len(out) == 0 || out[0] != val {
		require.Equalf(t, val, out, _notEqualHeader, key, out[0], val)

		return false
	}

	return true
}

// Headers assert value of the given header key:vak in the htt.Response param.
func Headers(t *testing.T, resp *http.Response, kv ...[2]string) bool {
	t.Helper()

	for i := range kv {
		k, v := kv[i][0], kv[i][1]
		if !Header(t, k, v, resp) {
			return false
		}
	}

	// if len(resp.Header) != len()

	return true
}

// Headers assert value of the given header key:vak in the htt.Response param.
func HeadersExact(t *testing.T, resp *http.Response, kv ...[2]string) bool {
	if !Headers(t, resp, kv...) {
		return false
	}

	require.Equalf(t, len(kv), len(resp.Header),
		"headers differs (exp != current): %+v != %+v",
		kv, resp.Header)

	return true
}

// DeleteAndTestAPI run a DELETE request before running the test handler.
func DeleteAndTestAPI(t *testing.T, url string, handler HandlerForTest) {
	t.Helper()

	resp := deleteAPI(t, url)
	defer resp.Body.Close()

	handler(t, resp)
}

// RequestAndTestAPI request an API then run the test handler.
func RequestAndTestAPI(t *testing.T, url string, handler HandlerForTest) {
	t.Helper()

	resp := requestAPI(t, url)
	defer resp.Body.Close()

	handler(t, resp)
}

// PushAndTestAPI post to an API then run the test handler.
// The sub method try to send an `application/json` encoded content
func PushAndTestAPI(t *testing.T, path string, content []byte, handler HandlerForTest, headers ...[2]string) {
	t.Helper()

	resp := pushAPI(t, path, content, headers...)
	defer resp.Body.Close()

	handler(t, resp)
}

// FetchBody return the response body.
func FetchBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	return fetchBody(t, resp)
}

func fetchBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	tmp, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error fetching the body response : %s", err.Error())
	}
	defer resp.Body.Close()

	return string(tmp)
}

func prepReq(t *testing.T, url, method string) *http.Response {
	t.Helper()

	var (
		client   = &http.Client{}
		req, err = http.NewRequest(method, url, http.NoBody)
	)

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
	t.Helper()

	return prepReq(t, url, "GET")
}

func deleteAPI(t *testing.T, url string) *http.Response {
	t.Helper()

	return prepReq(t, url, "DELETE")
}

func pushAPI(t *testing.T, url string, content []byte, headers ...[2]string) *http.Response {
	t.Helper()

	var (
		client   = &http.Client{}
		ctx, cl  = context.WithTimeout(context.Background(), time.Second*_ttlTest)
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(content))
	)

	defer cl()

	if err != nil {
		t.Fatalf("can't post the new request : %s", err.Error())
	}

	if len(headers) == 0 {
		req.Header.Set("Content-Type", "application/json")
	} else {
		for i := range headers {
			req.Header.Set(headers[i][0], headers[i][1])
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	return resp
}
