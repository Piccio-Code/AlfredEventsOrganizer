package database

import (
	"context"
	"github.com/Piccio-Code/AlfredEventsOranizer/backend/internal/waha"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type GroupModel struct {
	DB *pgxpool.Pool
}

type GroupDB struct {
	ID             string    `json:"ID,omitempty"`
	Title          string    `json:"title,omitempty" json:"title,omitempty"`
	WhatsappChatId string    `json:"whatsappChatId,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

type GroupRequest struct {
	Title          string `json:"title,omitempty"`
	WhatsappChatId string `json:"whatsappChatId,omitempty"`
}

func (m *GroupModel) Create(ctx context.Context, newGroup GroupRequest, participants []waha.GroupParticipant) (createGroup GroupDB, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	tx, err := m.DB.Begin(ctx)

	defer tx.Rollback(ctx)

	if err != nil {
		return GroupDB{}, err
	}

	stmt := `INSERT INTO groups(title, whatsapp_chat_id) 
				values ($1, $2) 
				RETURNING id, title, whatsapp_chat_id, created_at`

	err = tx.QueryRow(ctx, stmt, newGroup.Title, newGroup.WhatsappChatId).Scan(&createGroup.ID, &createGroup.Title, &createGroup.WhatsappChatId, &createGroup.CreatedAt)

	if err != nil {
		return GroupDB{}, err
	}

	for _, participant := range participants {
		err := m.AddParticipants(tx, ctx, createGroup, participant)

		if err != nil {
			return GroupDB{}, err
		}
	}

	err = tx.Commit(ctx)

	if err != nil {
		return GroupDB{}, err
	}

	return createGroup, nil
}

func (m *GroupModel) AddParticipants(tx pgx.Tx, ctx context.Context, group GroupDB, participant waha.GroupParticipant) error {

	stmt := `INSERT INTO group_members(group_id, name, whatsapp_user_id)  
			 VALUES ($1, $2, $3)`

	args := []interface{}{group.ID, participant.Name, participant.Id}

	_, err := tx.Exec(ctx, stmt, args...)

	return err
}
