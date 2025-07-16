package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

func MockInput(inputs ...string) func() {
	r, w, _ := os.Pipe()
	origStdin := os.Stdin
	os.Stdin = r
	go func() {
		for _, input := range inputs {
			if _, err := w.WriteString(input + "\n"); err != nil {
				fmt.Fprintf(os.Stderr, "failed to write mock input: %v\n", err)
			}
		}
		w.Close()
	}()
	return func() {
		os.Stdin = origStdin
	}
}

func newMockClient(statusCode int, body string) *resty.Client {
	client := resty.New()
	client.SetTransport(&mockRoundTripper{
		statusCode: statusCode,
		body:       body,
	})
	return client
}

type mockRoundTripper struct {
	statusCode int
	body       string
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.body)),
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")
	return resp, nil
}

type errorRoundTripper struct{}

func (e *errorRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("request failed")
}

func CaptureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	if _, err := buf.ReadFrom(r); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read captured output: %v\n", err)
	}
	os.Stdout = stdout
	return buf.String()
}
