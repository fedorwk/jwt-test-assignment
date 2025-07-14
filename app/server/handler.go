package server

import (
	"medods-auth/service/auth"
	"medods-auth/token"
	"medods-auth/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RefreshRequest struct {
	UserID       uuid.UUID `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type LogoutRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	AccessToken string    `json:"access_token"`
}

type MeRequest struct {
	AccessToken string `json:"access_token"`
}

func newGenerateHandler(authservice *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		guid := c.Query("guid")
		if guid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing guid parameter"})
			return
		}

		userID, err := uuid.Parse(guid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guid format"})
			return
		}

		u := user.User{
			Id:        userID,
			UserAgent: c.Request.UserAgent(),
		}

		pair, err := authservice.GenerateTokens(u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  string(*pair.Access),
			"refresh_token": string(*pair.Refresh),
		})
	}
}

func newRefreshHandler(authservice *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		u := user.User{
			Id:        req.UserID,
			UserAgent: c.Request.UserAgent(),
		}

		accessTok := token.EncodedToken(req.AccessToken)
		refreshTok := token.EncodedToken(req.RefreshToken)
		pair := auth.TokenPair{
			Access:  &accessTok,
			Refresh: &refreshTok,
		}

		newPair, err := authservice.Refresh(u, pair)
		if err != nil {
			status := http.StatusInternalServerError
			if err == auth.ErrUserAgentChanged ||
				err == auth.ErrUserIDMissmatch ||
				err == auth.ErrTokenExpired ||
				err == auth.ErrBlackListedToken {
				status = http.StatusUnauthorized
			}
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  string(*newPair.Access),
			"refresh_token": string(*newPair.Refresh),
		})
	}
}

func newMeHandler(authSvc *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		tokenStr := token.EncodedToken(req.AccessToken)
		id, err := authSvc.ExtractUserID(&tokenStr)
		if err != nil {
			status := http.StatusUnauthorized
			if err == auth.ErrTokenExpired ||
				err == auth.ErrBlackListedToken {
				status = http.StatusForbidden
			}
			c.JSON(status, gin.H{"error": "invalid token: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_id": id})
	}
}

func newLogoutHandler(authservice *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LogoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		u := user.User{
			Id:        req.UserID,
			UserAgent: c.Request.UserAgent(),
		}

		if err := authservice.RevokeTokens(u, token.EncodedToken(req.AccessToken)); err != nil {
			status := http.StatusInternalServerError
			if err == auth.ErrUserAgentChanged ||
				err == auth.ErrUserIDMissmatch ||
				err == auth.ErrTokenExpired {
				status = http.StatusUnauthorized
			}
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}
