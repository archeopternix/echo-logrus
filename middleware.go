// Package middleware provides echo request and response output log
package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo"
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
			fmt.Println("STEP")
			log.WithFields(logrus.Fields{
				"status":   res.Status,
				"method":   req.Method,
				"id":       id,
				"realip":   c.RealIP(),
				"duration": stop.Sub(start).String(),
			})
			log.Infof("*** %s %s [%v] %s %-7s %s %3d %s %s %13v %s %s",
				id,
				c.RealIP(),
				stop.Format(time.RFC3339),
				req.Host,
				req.Method,
				req.RequestURI,
				res.Status,
				reqSize,
				strconv.FormatInt(res.Size, 10),
				stop.Sub(start).String(),
				req.Referer(),
				req.UserAgent(),
			)
			return err
		}
	}
}
