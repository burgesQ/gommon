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
func Body(t *testing.T, expected string, resp *http.Response) {
	t.Helper()
	require.Equal(t, expected, fetchBody(t, resp))
}

func BodyContains(t *testing.T, expected string, resp *http.Response) {
	t.Helper()
	require.Contains(t, fetchBody(t, resp), expected)
}

// BodyDiffere fetch and assert that the body of the http.Response differ than expected
func BodyDiffere(t *testing.T, expected string, resp *http.Response) {
	t.Helper()
	require.NotEqual(t, expected, fetchBody(t, resp))
}

// StatusCode assert the status code of the response
func StatusCode(t *testing.T, expected int, resp *http.Response) {
	t.Helper()
	require.Equal(t, expected, resp.StatusCode)
}

// Header assert value of the given header key:vak in the htt.Response param
func Header(t *testing.T, key, val string, resp *http.Response) bool {
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
		if !Header(t, k, v, resp) {
			return false
		}
	}

	return true
}

func DeleteAndTestAPI(t *testing.T, url string, handler HandlerForTest) {
	t.Helper()

	resp := deleteAPI(t, url)
	defer resp.Body.Close()

	handler(t, resp)
}

// RequestAndTestAPI request an API then run the test handler
func RequestAndTestAPI(t *testing.T, url string, handler HandlerForTest) {
	t.Helper()

	resp := requestAPI(t, url)
	defer resp.Body.Close()

	handler(t, resp)
}

// PushAndTestAPI post to an API then run the test handler
// The sub method try to send an `application/json` encoded content
func PushAndTestAPI(t *testing.T, path string, content []byte, handler HandlerForTest, headers ...[2]string) {
	t.Helper()

	resp := pushAPI(t, path, content, headers...)
	defer resp.Body.Close()

	handler(t, resp)
}

func FetchBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	return fetchBody(t, resp)
}

func fetchBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	tmp, err := ioutil.ReadAll(resp.Body)
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
		ctx, cl  = context.WithTimeout(context.Background(), time.Second*_ttlTest)
		req, err = http.NewRequestWithContext(ctx, method, url, http.NoBody)
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
