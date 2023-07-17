package guard

import (
	syscontext "context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/pkg/errors"

	"github.com/yukels/ddos-guard/config"
	"github.com/yukels/util/collection"
	"github.com/yukels/util/context"
	"github.com/yukels/util/http/server"
	"github.com/yukels/util/log"
	"github.com/yukels/util/prometheus"
)

type Proxy struct {
	proxy       *goproxy.ProxyHttpServer
	guard       *Guard
	userService *UserService
	collector   *Collector
	config      *config.ProxyConfig
}

const (
	requestTimestamp = "requestStartTimestamp"
)

var (
	statusTeapot         = strconv.Itoa(http.StatusTeapot)
	statusTooManyRequest = strconv.Itoa(http.StatusTooManyRequests)
)

func NewDdosProxy(ctx context.Context) (*Proxy, error) {
	collector := NewCollector(ctx)

	guard, err := NewGuard(ctx, collector)
	if err != nil {
		return nil, err
	}
	userService, err := NewUserService(ctx, &config.Configs.DdosGuardConfig.UserService)
	if err != nil {
		return nil, err
	}

	p := &Proxy{
		guard:       guard,
		userService: userService,
		collector:   collector,
		config:      &config.Configs.DdosGuardConfig.Proxy,
	}

	if err = p.config.ReadFromEnv(ctx); err != nil {
		return nil, err
	}

	if err := p.init(ctx); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Proxy) Run(ctx context.Context) error {
	if err := p.userService.Run(ctx); err != nil {
		return err
	}
	if err := p.guard.Run(ctx); err != nil {
		return err
	}

	config := server.Config{
		Hostname: "",
		Port:     p.config.PortIn,
	}
	s := server.New(ctx, &config, p.proxy)
	server.ListenAndServe(ctx, s)

	return nil
}

func (p *Proxy) init(ctx context.Context) error {
	byPass := false
	byPassFromEnv := os.Getenv("BYPASS")
	if strings.ToLower(byPassFromEnv) == "true" {
		log.Log(ctx).Info("'BYPASS' mode ON")
		byPass = true
	}

	p.proxy = goproxy.NewProxyHttpServer()
	p.proxy.Tr.DialContext = func(ctx syscontext.Context, network, addr string) (c net.Conn, err error) {
		c, err = net.Dial(network, addr)
		if c, ok := c.(*net.TCPConn); err == nil && ok {
			c.SetKeepAlive(true)
		}
		return
	}
	p.proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Host == "" {
			fmt.Fprintln(w, "Cannot handle requests without Host header, e.g., HTTP 1.0")
			return
		}
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d", p.config.HostOut, p.config.PortOut)

		p.proxy.ServeHTTP(w, req)
	})

	if err := p.routes(ctx); err != nil {
		return err
	}

	minRetryAfter := p.config.RetryAfter / 2
	maxRetryAfter := p.config.RetryAfter

	p.proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(`.*`))).DoFunc(
		func(req *http.Request, proxyctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			// Stop guard - just forward the request
			if byPass {
				return req, nil
			}

			ctx := context.New(proxyctx.Req.Context())
			user := p.userService.UserFromRequest(ctx, req)
			if user == "" || collection.Contains(p.userService.GetWhiteListUsers(ctx), user) {
				return req, nil
			}

			if collection.Contains(p.userService.GetBlockedUsers(ctx), user) {
				resp := goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusTeapot, "Your user is blocked")
				p.collector.IncUserBlock(ctx, statusTeapot, user)
				return req, resp
			}

			proxyctx.Req = proxyctx.Req.WithContext(ctx.WithUser(user).WithValue(requestTimestamp, time.Now().UTC()))
			if p.guard.ShouldBlockUser(ctx, user) {
				resp := goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusTooManyRequests,
					"Too many request by user")
				// Use rand to spread incoming DDOS calls
				resp.Header.Add("Retry-After", strconv.Itoa(rand.Intn(maxRetryAfter-minRetryAfter)+minRetryAfter))
				p.collector.IncUserBlock(ctx, statusTooManyRequest, user)
				return req, resp
			}
			return req, nil
		})
	p.proxy.OnResponse(goproxy.UrlMatches(regexp.MustCompile(`.*`))).DoFunc(
		func(resp *http.Response, proxyctx *goproxy.ProxyCtx) *http.Response {
			if byPass {
				return resp
			}

			ctx := context.New(proxyctx.Req.Context())
			user := ctx.User()
			if user == nil || *user == "" {
				return resp
			}
			ts := ctx.Value(requestTimestamp)
			duration := time.Now().UTC().Sub(ts.(time.Time))
			p.guard.RequestCompleteUser(ctx, *user, duration)
			return resp
		})
	return nil
}

func (p *Proxy) routes(ctx context.Context) error {
	healthPath := p.config.HealthPath
	if healthPath != "" {
		p.proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("^%s$", healthPath)))).DoFunc(
			func(req *http.Request, proxyctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				return req, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusOK, "healthy")
			})
	}
	metricsPath := p.config.MetricsPath
	if metricsPath != "" {
		merticHandler, err := prometheus.NewHandler(ctx)
		if err != nil {
			return errors.Wrap(err, "Can't initialize metrics handler")
		}
		p.proxy.OnRequest(goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("^%s$", metricsPath)))).DoFunc(
			func(req *http.Request, proxyctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
				recorder := httptest.NewRecorder()
				merticHandler.Metrics(recorder, req)
				return req, recorder.Result()
			})
	}
	return nil
}
