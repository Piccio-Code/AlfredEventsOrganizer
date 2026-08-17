package server

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/Piccio-Code/AlfredEventsOranizer/backend/internal/waha"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
)

func (s *Server) AuthWebMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if gin.Mode() != "release" {
			user, err := s.authClient.GetUser(context.Background(), os.Getenv("UID_TEST_USER"))

			if err != nil {
				s.infoLog.Println(err)
				s.Fail(c, http.StatusUnauthorized, "Unauthorized User")
				c.Abort()
				return
			}

			c.Set(UserKey, user)
			c.Next()
			return
		}

		idToken := c.GetHeader("Authorization")
		idToken, ok := strings.CutPrefix(idToken, "Bearer ")

		if !ok || idToken == "" {
			s.Fail(c, http.StatusUnauthorized, "Unauthorized User")
			c.Abort()
			return
		}

		token, err := s.authClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			s.infoLog.Println(err)
			s.Fail(c, http.StatusUnauthorized, "Unauthorized User")
			c.Abort()
			return
		}

		user, err := s.authClient.GetUser(context.Background(), token.UID)

		if err != nil {
			s.infoLog.Println(err)
			s.Fail(c, http.StatusUnauthorized, "Unauthorized User")
			c.Abort()
			return
		}

		c.Set(UserKey, user)
		c.Next()
	}
}

func (s *Server) WahaSessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := s.wahaClient.GetWahaSession()

		if err != nil {
			s.infoLog.Println(err)
			s.Fail(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			c.Abort()
			return
		}

		if session.Status == waha.StatusWorking {
			c.Next()
			return
		}

		if session.Status == waha.StatusScanQRCode {
			s.Fail(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			c.Abort()
			return
		}

		err = s.wahaClient.RestartSession()

		if err != nil {
			s.Fail(c, http.StatusBadRequest, fmt.Sprintf("the session has a problem Status: %v", session.Status))
			c.Abort()
			return
		}
	}
}

func (s *Server) WahaHMACAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		Hmac := strings.TrimSpace(c.GetHeader("X-Webhook-Hmac"))

		body, err := io.ReadAll(c.Request.Body)

		if err != nil {
			s.infoLog.Println(err)
			s.Fail(c, http.StatusUnauthorized, "Unauthorized User")
			c.Abort()
			return
		}

		messageMAC, err := hex.DecodeString(Hmac)

		if err != nil {
			s.infoLog.Println(err)
			s.Fail(c, http.StatusUnauthorized, "Unauthorized User")
			c.Abort()
			return
		}

		isValid := ValidMAC(body, messageMAC, []byte(os.Getenv("WAHA_WEBHOOK_HMAC_KEY")))

		if !isValid {
			s.Fail(c, http.StatusUnauthorized, "Unauthorized User")
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		c.Next()
	}
}
