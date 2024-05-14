package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	coreParams "github.com/smalls0098/proxy-download/params"
	pkgApp "github.com/smalls0098/proxy-download/pkg/app"
	pkgHttp "github.com/smalls0098/proxy-download/pkg/app/server/http"
	"golang.org/x/time/rate"
	"io"
	"log"
	"net"
	"net/http"
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

	s.Any("/proxy", handleProxy)

	app := pkgApp.New(
		pkgApp.WithServer(s),
		pkgApp.WithName("proxy"),
	)
	log.Printf("running: [http://127.0.0.1:%d]", p)
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}

var hc = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("stopped after 5 redirects")
		}
		return nil
	},
}

var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

func handleProxy(ctx *gin.Context) {
	f := ctx.Query("f")
	ps := ctx.Query("p")
	if len(ps) == 0 {
		ctx.String(http.StatusBadRequest, "参数不能为空")
		ctx.Abort()
		return
	}
	params, err := coreParams.DecParams(ps, key)
	if err != nil {
		handleErr(ctx, err)
		return
	}
	reqUrl, err := url.Parse(params.Url)
	if err != nil {
		handleErr(ctx, err)
		return
	}
	log.Println(ctx.ClientIP() + " " + ctx.Request.Method + " " + reqUrl.String() + " " + ctx.Request.Proto + " " + ctx.Request.UserAgent())

	req, err := http.NewRequest(ctx.Request.Method, reqUrl.String(), http.NoBody)
	if err != nil {
		handleErr(ctx, err)
		return
	}
	req.Header = ctx.Request.Header.Clone()
	req.Header.Del("Host")
	req.Header.Del("Referer")
	req.Header.Del("Origin")

	res, err := hc.Do(req)
	if err != nil {
		handleErr(ctx, err)
		return
	}
	defer res.Body.Close()

	// 重定向
	if res.StatusCode >= 300 && res.StatusCode < 400 {
		if location, _ := res.Location(); location != nil {
			loc := location.String()
			loc, err = coreParams.Gen(loc, f, params, key)
			if err != nil {
				handleErr(ctx, err)
				return
			}
			ctx.Redirect(res.StatusCode, loc)
			return
		}
	}

	// 设置headers
	resHr := res.Header.Clone()
	for _, k := range hopHeaders {
		resHr.Del(k)
	}
	for k, v := range resHr {
		for _, vv := range v {
			ctx.Writer.Header().Add(k, vv)
		}
	}

	ctx.Status(res.StatusCode)
	if ctx.Request.Method == http.MethodHead {
		return
	}

	if params.Tag == 1 {
		// 无限制
		_, err = io.Copy(ctx.Writer, res.Body)
	} else {
		// 默认限制，每秒1M
		_, err = io.Copy(ctx.Writer, &rateLimitedReader{
			ctx:     ctx,
			reader:  res.Body,
			limiter: rate.NewLimiter(rate.Limit(1024*1024), 1024*1024),
		})
	}
	if err != nil {
		handleErr(ctx, err)
		return
	}
}

func handleErr(ctx *gin.Context, err error) {
	ctx.String(http.StatusBadGateway, "出现异常: %s", err.Error())
	ctx.Abort()
}

type rateLimitedReader struct {
	ctx     context.Context
	reader  io.Reader
	limiter *rate.Limiter
}

func (r *rateLimitedReader) Read(p []byte) (n int, err error) {
	err = r.limiter.WaitN(r.ctx, len(p))
	if err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}
