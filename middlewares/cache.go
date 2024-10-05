package middlewares

import (
	"bytes"
	"context"
	"regexp"
	"time"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/redis/go-redis/v9"
)

import (
	"github.com/huangcheng/icalhub/config"
	"github.com/huangcheng/icalhub/utils"
)

var ctx = context.Background()

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func CacheMiddleware(config config.Config) gin.HandlerFunc {
	redisClient := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    config.RedisHost + ":" + config.RedisPort,
		DB:      config.RedisDB,
	})

	ttl := time.Duration(config.CacheTTL) * time.Second

	response := new(utils.Response)

	return func(c *gin.Context) {
		key := c.Request.RequestURI

		match := regexp.MustCompile(`^/(holidays|movies)/`).MatchString(key)

		if value, err := redisClient.Get(ctx, key).Result(); err == nil && len(value) > 0 && match {
			response.Calendar(c, []byte(value))

			c.AbortWithStatus(200)
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		if c.Writer.Status() == 200 && match {
			body := w.body.String()

			if len(body) > 0 {
				redisClient.SetEx(ctx, key, body, ttl)
			}
		}
	}
}
