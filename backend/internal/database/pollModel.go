package database

import (
	"time"

	"github.com/Piccio-Code/AlfredEventsOranizer/backend/internal/waha"
)

type PollRequest struct {
	TemplateID     *string             `json:"template_id,omitempty"`
	Name           string              `json:"name,omitempty"`
	Title          string              `json:"title,omitempty"`
	GroupID        string              `json:"groupId,omitempty"`
	MultipleChoice bool                `json:"multipleChoice,omitempty"`
	Options        []PollOptionRequest `json:"options,omitempty"`
	ExpiresAt      time.Time           `json:"expiresAt,omitempty"`
}

type PollOptionRequest struct {
	Label                string `json:"label,omitempty"`
	Position             int    `json:"position,omitempty"`
	MotivationNeeded     *bool  `json:"motivationNeeded,omitempty"`
	CongratulationNeeded *bool  `json:"congratulationNeeded,omitempty"`
	SpecificationNeeded  *bool  `json:"specificationNeeded,omitempty"`
}

func (r PollRequest) ToWahaPoll() waha.NewPollRequest {
	options := make([]string, 0, len(r.Options))
	for _, option := range r.Options {
		options = append(options, option.Label)
	}

	return waha.NewPollRequest{
		ChatId: r.GroupID,
		Poll: waha.PollPayload{
			Name:            r.Title,
			Options:         options,
			MultipleAnswers: r.MultipleChoice,
		},
	}
}

type PollDB struct {
	ID             string
	TemplateID     string
	WhatsAppPollID string
	Title          string
	MultipleChoice bool
	CreatedAt      time.Time
	ExpiresAt      time.Time
	GroupID        string
	Status         string
}
