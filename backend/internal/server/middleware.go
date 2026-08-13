package server

import (
	"bytes"
	"context"
	"encoding/hex"
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
		session, err := s.GetWahaSession()

		if err != nil {
			s.Fail(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			c.Abort()
			return
		}

		if session.Status == "WORKING" {
			c.Next()
			return
		}

		if session.Status == "STOPPED" || session.Status == "FAILED" {
			err := s.RestartSession()

			if err != nil {
				s.Fail(c, http.StatusBadRequest, "The session is off")
				c.Abort()
				return
			}
		}

		s.Fail(c, http.StatusBadRequest, "The session is off check telegram for the login code")
		c.Abort()
	}
}

func (s *Server) WahaHMACAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		Hmac := strings.TrimSpace(c.GetHeader("X-Webhook-Hmac"))

		s.infoLog.Println(Hmac)

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
