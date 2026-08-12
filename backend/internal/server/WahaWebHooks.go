package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SessionWebHook struct {
	Event   string `json:"event"`
	Session string `json:"session"`
	Me      struct {
		Id       string `json:"id"`
		PushName string `json:"pushName"`
	} `json:"me"`
	Payload struct {
		Status   string `json:"status"`
		Statuses []struct {
			Status    string `json:"status"`
			Timestamp int64  `json:"timestamp"`
		} `json:"statuses"`
		Data interface{} `json:"data"`
	} `json:"payload"`
	Engine      string `json:"engine"`
	Environment struct {
		Version string `json:"version"`
		Engine  string `json:"engine"`
		Tier    string `json:"tier"`
	} `json:"environment"`
}

func (s *Server) SessionWebHookHandler(c *gin.Context) {
	var session SessionWebHook
	if err := c.ShouldBindJSON(&session); err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "Error Biding The Session")
		return
	}

	s.infoLog.Println(session)

	if session.Payload.Status == "SCAN_QR_CODE" {
		code, err := s.getSessionCode()

		if err != nil {
			return
		}

		s.infoLog.Println(code)
	}
}
