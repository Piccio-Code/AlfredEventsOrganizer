package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"net/http"
	"os"
	"time"
)

type SessionWebHook struct {
	Event          string    `json:"event"`
	Timestamp      time.Time `json:"timestamp"`
	SessionId      string    `json:"sessionId"`
	IdempotencyKey string    `json:"idempotencyKey"`
	DeliveryId     string    `json:"deliveryId"`
	Data           struct {
		SessionId string `json:"sessionId"`
		Status    string `json:"status"`
	} `json:"data"`
}

func (s *Server) RegisterWahaRoutes(r *gin.RouterGroup) {
	webhooks := r.Group("/webhooks")
	{
		webhooks.POST("/session", s.SessionWebHookHandler)
	}

	client := r.Group("")
	{
		client.Use(s.WahaSessionCheck())

		client.GET("/test",
			func(context *gin.Context) {
				s.Ok(context, nil, nil)
				return
			})
	}
}

func (s *Server) SessionWebHookHandler(c *gin.Context) {
	var session SessionWebHook

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

	code, err := s.GetSessionCode()
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
