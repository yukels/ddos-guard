package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yukels/util/context"
)

const (
	port = 18888
)

func TestPost(t *testing.T) {
	values := url.Values{"key": []string{"val"}}
	tests := []struct {
		config      *Config
		timeout     time.Duration
		expected    int
		expectedErr bool
	}{
		// Good
		{&Config{Hostname: "", Port: port, ReadTimeout: 1}, 0, http.StatusOK, false},
		// MaxBytes OK
		{&Config{Hostname: "", Port: port, MaxBytes: 7}, 0, http.StatusOK, false},
		// MaxBytes failed
		{&Config{Hostname: "", Port: port, MaxBytes: 6}, 0, http.StatusBadRequest, false},
		// WriteTimeout failed
		{&Config{Hostname: "", Port: port, WriteTimeout: 1}, 2, http.StatusBadRequest, true},
	}
	for idx, tst := range tests {
		ctx := context.Background()
		s := New(ctx, tst.config, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			for k, v := range values {
				value := r.PostFormValue(k)
				if v[0] != value {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			select {
			case <-time.After(tst.timeout * time.Second):
				w.WriteHeader(http.StatusOK)
			case <-ctx.Done():
				w.WriteHeader(http.StatusRequestTimeout)
			}
		}))

		go func() {
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				assert.NoErrorf(t, err, "[%d] Error on server listening ...", idx)
			}
		}()

		time.Sleep(time.Second)

		url := fmt.Sprintf("http://%s:%d", tst.config.Hostname, tst.config.Port)

		req, err := http.PostForm(url, values)
		if (err != nil) != tst.expectedErr {
			t.Errorf("[%d]: expected error status %v, got %v", idx, tst.expectedErr, (err != nil))
		}

		if req != nil {
			assert.Equalf(t, tst.expected, req.StatusCode, "[%d] Unexpected StatusCode", idx)
		}

		ctxShutdown, _ := context.WithTimeout(ctx, time.Second)
		s.Shutdown(ctxShutdown)
	}
}

type SlowReader struct {
	r       io.Reader
	timeout time.Duration
}

func NewSlowReader(data string, timeout time.Duration) io.Reader {
	return &SlowReader{r: strings.NewReader(data), timeout: timeout}
}

func (r *SlowReader) Read(buf []byte) (n int, err error) {
	time.Sleep(r.timeout)
	return r.r.Read(buf)
}

func TestTimeout(t *testing.T) {
	tests := []struct {
		config      *Config
		body        string
		timeout     time.Duration
		expected    int
		expectedErr bool
	}{
		// Good
		{&Config{Hostname: "", Port: port, ReadTimeout: 1}, "foo=bar", 0, http.StatusOK, false},
		// Read timeout
		{&Config{Hostname: "", Port: port, ReadTimeout: 1}, "foo=bar", 2 * time.Second, http.StatusRequestTimeout, false},
	}
	for idx, tst := range tests {
		ctx := context.Background()
		s := New(ctx, tst.config, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()

			if err != nil {
				if netErr, ok := err.(*net.OpError); ok && netErr.Timeout() {
					w.WriteHeader(http.StatusRequestTimeout)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
				return
			}
			if string(b) != tst.body {
				w.WriteHeader(http.StatusBadRequest)
			}
		}))

		go func() {
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				assert.NoErrorf(t, err, "[%d] Error on server listening ...", idx)
			}
		}()

		time.Sleep(time.Second)

		url := fmt.Sprintf("http://%s:%d", tst.config.Hostname, tst.config.Port)

		req, err := http.NewRequest("GET", url, NewSlowReader(tst.body, tst.timeout))
		assert.NoErrorf(t, err, "[%d] Error on NewRequest", idx)

		resp, err := http.DefaultClient.Do(req)
		if (err != nil) != tst.expectedErr {
			t.Errorf("[%d]: expected error status %v, got %v", idx, tst.expectedErr, (err != nil))
		}

		if resp != nil {
			defer resp.Body.Close()
			assert.Equalf(t, tst.expected, resp.StatusCode, "[%d] Unexpected StatusCode", idx)
		}

		ctxShutdown, _ := context.WithTimeout(ctx, time.Second)
		s.Shutdown(ctxShutdown)
	}
}
