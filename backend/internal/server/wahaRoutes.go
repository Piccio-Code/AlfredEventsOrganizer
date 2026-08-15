package server

import (
	"fmt"
	"github.com/Piccio-Code/AlfredEventsOranizer/backend/internal/waha"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"net/http"
	"os"
)

func (s *Server) RegisterWahaRoutes(r *gin.RouterGroup) {
	webhooks := r.Group("/webhooks")
	{
		webhooks.POST("/session", s.SessionWebHookHandler)
	}

	client := r.Group("")
	{
		client.Use(s.WahaSessionCheck())

		client.GET("/groups-list", s.ListGroups)

	}
}

func (s *Server) SessionWebHookHandler(c *gin.Context) {
	var session waha.SessionWebHook

	if err := c.ShouldBindJSON(&session); err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusBadRequest, "error binding the session")
		return
	}

	s.infoLog.Println("Session status:", session.Data.Status)

	if session.Data.Status != "qr_ready" {
		s.Ok(c, nil, nil)
		return
	}

	code, err := s.wahaClient.GetSessionCode()
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusInternalServerError, "error getting the WhatsApp pairing code")
		return
	}

	if code.Status != "qr_ready" {
		s.infoLog.Println("Unexpected pairing code status:", code.Status)
		s.Fail(c, http.StatusInternalServerError, "WhatsApp pairing code is not ready")
		return
	}

	_, err = s.alfredTelegram.SendMessage(
		c.Request.Context(),
		&bot.SendMessageParams{
			ChatID: os.Getenv("telegram_chat_id"),
			Text: fmt.Sprintf(
				"🎩 <b>Autenticazione WhatsApp</b>\n\n"+
					"Signore, il codice richiesto è il seguente:\n\n"+
					"<code>%s</code>\n\n"+
					"🧐 La procedura può ora essere completata.",
				code.PairingCode,
			),
			ParseMode: "HTML",
		},
	)

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusInternalServerError, "error sending the pairing code to Telegram")
		return
	}

	s.Ok(c, nil, nil)
}

func (s *Server) ListGroups(c *gin.Context) {
	groups, err := s.wahaClient.GetGroupsList()

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusInternalServerError, "error getting the groups")
		return
	}

	s.Ok(c, envelop{"groups": groups}, nil)
}
