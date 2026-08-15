package database

import (
	"context"
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

func (m *GroupModel) Create(ctx context.Context, newGroup GroupRequest) (createGroup GroupDB, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

	defer cancel()

	stmt := `INSERT INTO groups(title, whatsapp_chat_id) 
				values ($1, $2) 
				RETURNING id, title, whatsapp_chat_id, created_at`

	err = m.DB.QueryRow(ctx, stmt, newGroup.Title, newGroup.WhatsappChatId).Scan(&createGroup.ID, &createGroup.Title, &createGroup.WhatsappChatId, &createGroup.CreatedAt)

	if err != nil {
		return GroupDB{}, err
	}

	return createGroup, nil
}
