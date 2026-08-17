package server

import (
	"fmt"
	"github.com/Piccio-Code/AlfredEventsOranizer/backend/internal/waha"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"net/http"
	"os"
)

const sessionWorkingMessage = "🎩 <b>Sessione WhatsApp Operativa</b>\n\n" +
	"Signore, la sessione di WhatsApp è stata avviata correttamente e funziona senza alcun problema.\n\n" +
	"📡 Connessione: <b>stabile</b>\n" +
	"🔐 Autenticazione: <b>completata</b>\n" +
	"⚙️ Stato: <b>operativo</b>\n\n" +
	"🧐 Può procedere tranquillamente: il sistema è pronto e pienamente operativo."

const pairingCodeMessage = "🎩 <b>Autenticazione WhatsApp</b>\n\n" +
	"Signore, il codice richiesto è il seguente:\n\n" +
	"<code>%s</code>\n\n" +
	"🧐 La procedura può ora essere completata."

func (s *Server) RegisterWahaRoutes(r *gin.RouterGroup) {
	webhooks := r.Group("/webhooks")
	{
		webhooks.Use(s.WahaHMACAuth())

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

	s.infoLog.Println("Session status:", session.Payload.Status)

	switch {
	case session.BecameWorking():
		s.notifyTelegram(c, sessionWorkingMessage)
	case session.NeedsPairing():
		s.notifyPairingCode(c)
	default:
		s.Ok(c, nil, nil)
	}
}

func (s *Server) notifyPairingCode(c *gin.Context) {
	code, err := s.wahaClient.GetSessionCode()
	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusInternalServerError, "error getting the WhatsApp pairing code")
		return
	}

	if code.Code == "" {
		s.infoLog.Println("WAHA returned an empty pairing code")
		s.Fail(c, http.StatusInternalServerError, "WhatsApp pairing code is not ready")
		return
	}

	s.notifyTelegram(c, fmt.Sprintf(pairingCodeMessage, code.Code))
}

func (s *Server) notifyTelegram(c *gin.Context, text string) {
	_, err := s.alfredTelegram.SendMessage(
		c.Request.Context(),
		&bot.SendMessageParams{
			ChatID:    os.Getenv("TELEGRAM_CHAT_ID"),
			Text:      text,
			ParseMode: "HTML",
		},
	)

	if err != nil {
		s.infoLog.Println(err)
		s.Fail(c, http.StatusInternalServerError, "error sending the message to Telegram")
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
