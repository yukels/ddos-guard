package server

import (
	syscontext "context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
)

type maxBytesHandler struct {
	h http.Handler
	n int64
}

func (h *maxBytesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.n)
	h.h.ServeHTTP(w, r)
}

// New create new http server instance according to the config settings
func New(ctx context.Context, config *Config, handler http.Handler) *http.Server {
	endpoint := fmt.Sprintf("%s:%d", config.Hostname, config.Port)
	s := http.Server{
		Addr:    endpoint,
		Handler: handler,
	}
	if config.MaxBytes != 0 {
		s.Handler = &maxBytesHandler{h: handler, n: config.MaxBytes}
	}
	if config.ReadHeaderTimeout != 0 {
		s.ReadHeaderTimeout = time.Duration(config.ReadHeaderTimeout) * time.Second
	}
	if config.ReadTimeout != 0 {
		s.ReadTimeout = time.Duration(config.ReadTimeout) * time.Second
	}
	if config.WriteTimeout != 0 {
		s.WriteTimeout = time.Duration(config.WriteTimeout) * time.Second
	}
	if config.MaxHeaderBytes != 0 {
		s.MaxHeaderBytes = config.MaxHeaderBytes
	}

	log.Log(ctx).Infof("Listening on %s", endpoint)
	toLog, _ := json.MarshalIndent(config, "", "    ")
	log.Log(ctx).Infof("Config %s", string(toLog))
	return &s
}

// ListenAndServe listen to requests and gracefully shuts down the server on process interrupt
func ListenAndServe(ctx context.Context, s *http.Server) {
	listenErrChan := make(chan error)
	go func() {
		listenErrChan <- s.ListenAndServe()
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case <-sig:
		break
	case <-ctx.Done():
		break
	case listenErr := <-listenErrChan:
		if listenErr != nil && listenErr == http.ErrServerClosed {
			break
		}
		log.Log(ctx).WithError(errors.Errorf("Serve error: %s", listenErr)).Fatal("Error on REST listen. Server is stopped.")
	}

	shutdown := func(srv *http.Server, wg *sync.WaitGroup) {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err == syscontext.DeadlineExceeded {
			log.Log(ctx).WithError(err).Errorf("Force shutdown %s", srv.Addr)
		} else {
			log.Log(ctx).Debugf("Graceful shutdown %s", srv.Addr)
		}
		wg.Done()
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go shutdown(s, wg)
	wg.Wait()
}
