package server

import (
	"context"
	"github.com/gin-gonic/gin"
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

func (s *Server) WhatAppsSessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := s.GetWahaSession()

		if err != nil {
			s.Fail(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			c.Abort()
			return
		}

		if session.Status == "WORKING" {
			c.Next()
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
