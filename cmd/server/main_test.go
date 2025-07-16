package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/shekshuev/gophkeeper/internal/config"
	"github.com/stretchr/testify/require"
)

func Test_printBuildInfo(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printBuildInfo()

	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)

	output := buf.String()
	if !strings.Contains(output, "Build version") ||
		!strings.Contains(output, "Build date") ||
		!strings.Contains(output, "Build commit") {
		t.Errorf("Output missing expected strings: %s", output)
	}
}

func Test_Run_gracefulShutdown(t *testing.T) {
	_ = os.Setenv("SERVER_ADDRESS", "127.0.0.1:8089")
	_ = os.Setenv("DATABASE_DSN", "user=test dbname=test sslmode=disable")
	_ = os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	_ = os.Setenv("REFRESH_TOKEN_EXPIRES", "1h")
	_ = os.Setenv("ACCESS_TOKEN_SECRET", "access")
	_ = os.Setenv("REFRESH_TOKEN_SECRET", "refresh")

	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")

	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	err := cmd.Start()
	if err != nil {
		t.Fatalf("could not start subprocess: %v", err)
	}

	go func() {
		_ = cmd.Process.Signal(os.Interrupt)
	}()

	err = cmd.Wait()
	if err != nil && !strings.Contains(err.Error(), "signal: interrupt") {
		t.Errorf("subprocess exited with error: %v", err)
	}
}

func TestServer_HealthCheck(t *testing.T) {
	_ = os.Setenv("SERVER_ADDRESS", "127.0.0.1:8090")
	_ = os.Setenv("DATABASE_DSN", "user=test dbname=test sslmode=disable")
	_ = os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	_ = os.Setenv("REFRESH_TOKEN_EXPIRES", "1h")
	_ = os.Setenv("ACCESS_TOKEN_SECRET", "access")
	_ = os.Setenv("REFRESH_TOKEN_SECRET", "refresh")

	cfg := config.GetConfig()
	srv := NewServer(&cfg)

	go func() {
		_ = srv.ListenAndServe()
	}()
	defer srv.Close()

	time.Sleep(200 * time.Millisecond)

	resp, err := http.Get("http://" + cfg.ServerAddress + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, string(body), "ok")
}
