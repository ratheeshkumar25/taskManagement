package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/ratheeshkumar25/task-mgt/utility"
)

func RateLimitMiddleware(redisClient *redis.Client, maxRequests int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		userID := fmt.Sprintf("%v", c.GetString("userID"))
		if userID == "" || userID == "0" {
			userID = c.ClientIP() // fallback to IP
		}

		key := fmt.Sprintf("rate_limit:%v", userID)
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limiter internal error"})
			c.Abort()
			return
		}

		if count == 1 {
			_, err := redisClient.Expire(ctx, key, duration).Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limiter expiration error"})
				c.Abort()
				return
			}
		}

		if count > int64(maxRequests) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded. Try again later."})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AuthMiddleware checks for valid JWT in Authorization header and sets userID in context
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenStr, &utility.UserClaim{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*utility.UserClaim)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Attach UserID to request context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
