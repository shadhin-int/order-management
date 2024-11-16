package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"order-management/config"
	"order-management/models"
	"strings"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Message: "Unauthorized",
				Type:    "error",
				Code:    401,
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Message: "Unauthorized",
				Type:    "error",
				Code:    401,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
