// Package middleware provides echo request and response output log
package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/neko-neko/echo-logrus/v2/log"
	logrus "github.com/sirupsen/logrus"
)

// Logger returns a middleware that logs HTTP requests.
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			reqSize := req.Header.Get(echo.HeaderContentLength)
			if reqSize == "" {
				reqSize = "0"
			}
			log.InfoWithFields(logrus.Fields{
				"status":    res.Status,
				"method":    req.Method,
				"id":        id,
				"realip":    c.RealIP(),
				"duration":  stop.Sub(start).String(),
				"size":      strconv.FormatInt(res.Size, 10),
				"referrer":  req.Referer(),
				"host":      req.Host,
				"request":   req.RequestURI,
				"useragent": req.UserAgent(),
			})
			return err
		}
	}
}
