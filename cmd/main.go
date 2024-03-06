package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	coreParams "github.com/smalls0098/proxy-download/params"
	pkgApp "github.com/smalls0098/proxy-download/pkg/app"
	pkgHttp "github.com/smalls0098/proxy-download/pkg/server/http"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	key string
	p   int
)

func init() {
	flag.StringVar(&key, "key", "smalls0098", "default key smalls0098")
	flag.IntVar(&p, "p", 13822, "default port 13822")
	if v := os.Getenv("KEY"); len(v) > 0 {
		key = v
	}
	if v := os.Getenv("PORT"); len(v) > 0 {
		port, err := strconv.Atoi(v)
		if err == nil {
			p = port
		}
	}
}

func main() {
	// 执行命令行
	flag.Parse()

	server := &httputil.ReverseProxy{
		Rewrite:  nil,
		Director: func(req *http.Request) {},
		Transport: &retryRoundTripper{
			transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   20 * time.Second, // 请求超时
					KeepAlive: 20 * time.Second, // 检测连接是否存活
				}).DialContext,
				ForceAttemptHTTP2:     true,             // 强制尝试http2
				MaxIdleConns:          50,               // 最大空闲链接
				IdleConnTimeout:       30 * time.Second, // 空闲链接时间
				TLSHandshakeTimeout:   5 * time.Second,  // TLS握手超时时间
				ExpectContinueTimeout: 5 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableKeepAlives: true,
			},
			maxRetries: 2,
		},
		ErrorHandler: func(w http.ResponseWriter, req *http.Request, err error) {
			log.Printf("ErrorHandler: %+v, url: %s", err, req.URL)
		},
	}

	gin.SetMode(gin.ReleaseMode)
	s := pkgHttp.NewServer(
		gin.New(),
		pkgHttp.WithServerHost("0.0.0.0"),
		pkgHttp.WithServerPort(p),
		pkgHttp.WithServerTimeout(20*time.Second),
	)
	s.Use(gin.Recovery())
	s.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "not found")
	})
	s.NoMethod(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "not method")
	})
	s.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "nw smalls\nok")
	})

	s.Any("/proxy", proxy(server))

	app := pkgApp.New(
		pkgApp.WithServer(s),
		pkgApp.WithName("proxy"),
	)
	log.Printf("running: [http://127.0.0.1:%d]", p)
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}

type retryRoundTripper struct {
	transport  http.RoundTripper
	maxRetries int
}

func (r *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req != nil {
		req.Header.Del("X-Forwarded-For")
		req.Header.Del("Host")
		req.Header.Set("Host", req.URL.Host)
	}
	var redirectCount = 0
	var retryCount = 0
	var resp *http.Response
	var err error
	for {
		if redirectCount >= 3 {
			return nil, errors.New("stopped after 3 redirects")
		}
		if retryCount >= 2 {
			break
		}
		resp, err = r.transport.RoundTrip(req)
		// retry
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			retryCount++
			continue
		}
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			req.URL, err = resp.Location()
			if err != nil {
				return nil, err
			}
			redirectCount++
			continue
		}
		if 200 < resp.StatusCode || resp.StatusCode >= 400 {
			time.Sleep(500 * time.Millisecond)
			retryCount++
			continue
		}
		break
	}
	if err == nil && resp == nil {
		return nil, errors.New("出现异常: resp is nil")
	}
	return resp, err
}

func proxy(server *httputil.ReverseProxy) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ps := ctx.Query("p")
		if len(ps) == 0 {
			ctx.String(http.StatusBadRequest, "参数不能为空")
			return
		}
		params, err := coreParams.DecParams(ps, key)
		if err != nil {
			ctx.String(http.StatusBadGateway, "出现异常: %s", err.Error())
			return
		}

		reqUrl, err := url.Parse(params.Url)
		if err != nil {
			ctx.String(http.StatusBadGateway, "出现异常: %s", err.Error())
			return
		}

		// 请求体
		req := ctx.Request.Clone(ctx.Request.Context())
		req.URL = reqUrl

		server.ServeHTTP(ctx.Writer, req)
	}
}
