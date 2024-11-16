package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"order-management/config"
	"order-management/models"

	"time"
)

func Login(c *gin.Context) {
	var loginReq models.LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "Invalid request format",
			Type:    "error",
			Code:    400,
		})
		return
	}

	testUsername := config.AppConfig.TestAccess.TestUsername
	testPassword := config.AppConfig.TestAccess.TestPassword

	if loginReq.Username != testUsername || loginReq.Password != testPassword {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "The user credentials were incorrect",
			Type:    "error",
			Code:    400,
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginReq.Username,
		"exp":      time.Now().Add(config.AppConfig.JWT.ExpirationDur).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Message: "Error generating token",
			Type:    "error",
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		TokenType:    "Bearer",
		ExpiresIn:    43200,
		AccessToken:  tokenString,
		RefreshToken: "refresh_" + tokenString,
	})
}
