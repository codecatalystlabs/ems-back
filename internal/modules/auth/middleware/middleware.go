package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	platformauth "dispatch/internal/platform/auth"
	"dispatch/internal/platform/config"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	jwtMgr := platformauth.NewJWTManager(secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	return func(c *gin.Context) {
		auth := strings.TrimSpace(c.GetHeader("Authorization"))
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing bearer token"})
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimSpace(parts[1])
		claims, err := jwtMgr.ParseAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

func DeviceContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := strings.TrimSpace(c.GetHeader("X-Device-ID"))
		deviceName := strings.TrimSpace(c.GetHeader("X-Device-Name"))

		if deviceID == "" {
			deviceID = strings.TrimSpace(c.Query("device_id"))
		}
		if deviceName == "" {
			deviceName = strings.TrimSpace(c.Query("device_name"))
		}

		if deviceID == "" {
			deviceID = clientFingerprint(c)
		}
		if deviceName == "" {
			deviceName = inferDeviceName(c.Request.UserAgent())
		}

		c.Set("device_id", deviceID)
		c.Set("device_name", deviceName)
		c.Next()
	}
}

func clientFingerprint(c *gin.Context) string {
	ip := c.ClientIP()
	ua := c.Request.UserAgent()
	if ua == "" {
		ua = "unknown-agent"
	}
	return ip + "|" + ua
}

func inferDeviceName(userAgent string) string {
	ua := strings.ToLower(strings.TrimSpace(userAgent))
	switch {
	case strings.Contains(ua, "android"):
		return "Android Device"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ios"):
		return "iOS Device"
	case strings.Contains(ua, "windows"):
		return "Windows Device"
	case strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os"):
		return "Mac Device"
	case strings.Contains(ua, "linux"):
		return "Linux Device"
	default:
		return "Unknown Device"
	}
}
