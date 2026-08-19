package database

import (
	"context"
	"errors"
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
	Title          string    `json:"title,omitempty"`
	WhatsappChatId string    `json:"whatsappChatId,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

type GroupRequest struct {
	Title          string `json:"title,omitempty"`
	WhatsappChatId string `json:"whatsappChatId,omitempty"`
}

type Member struct {
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	WhatsappUserId string    `json:"whatsappUserId,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

type GroupResponse struct {
	ID             string    `json:"id,omitempty"`
	Title          string    `json:"title,omitempty"`
	WhatsappChatId string    `json:"whatsappChatId,omitempty"`
	Members        []Member  `json:"members,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
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

func (m *GroupModel) GetAll(ctx context.Context) ([]GroupResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `SELECT id, title, whatsapp_chat_id, created_at 
			 FROM groups 
			 ORDER BY created_at DESC`

	rows, err := m.DB.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []GroupResponse{}
	groupIndex := map[string]int{}

	for rows.Next() {
		var group GroupResponse

		err := rows.Scan(&group.ID, &group.Title, &group.WhatsappChatId, &group.CreatedAt)
		if err != nil {
			return nil, err
		}

		group.Members = []Member{}
		groupIndex[group.ID] = len(groups)
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows.Close()

	stmt = `SELECT id, group_id, name, whatsapp_user_id, created_at 
			FROM group_members`

	memberRows, err := m.DB.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer memberRows.Close()

	for memberRows.Next() {
		var member Member
		var groupID string

		err := memberRows.Scan(&member.ID, &groupID, &member.Name, &member.WhatsappUserId, &member.CreatedAt)
		if err != nil {
			return nil, err
		}

		idx, ok := groupIndex[groupID]
		if !ok {
			continue
		}

		groups[idx].Members = append(groups[idx].Members, member)
	}

	if err := memberRows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (m *GroupModel) GetByID(ctx context.Context, id string) (GroupResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `SELECT id, title, whatsapp_chat_id, created_at 
			 FROM groups 
			 WHERE id = $1`

	var group GroupResponse

	err := m.DB.QueryRow(ctx, stmt, id).Scan(&group.ID, &group.Title, &group.WhatsappChatId, &group.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return GroupResponse{}, NotFoundError
	}

	if err != nil {
		return GroupResponse{}, err
	}

	stmt = `SELECT id, name, whatsapp_user_id, created_at 
			FROM group_members 
			WHERE group_id = $1`

	rows, err := m.DB.Query(ctx, stmt, id)
	if err != nil {
		return GroupResponse{}, err
	}
	defer rows.Close()

	members := []Member{}

	for rows.Next() {
		var member Member

		err := rows.Scan(&member.ID, &member.Name, &member.WhatsappUserId, &member.CreatedAt)
		if err != nil {
			return GroupResponse{}, err
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return GroupResponse{}, err
	}

	group.Members = members

	return group, nil
}

func (m *GroupModel) Update(ctx context.Context, id string, title string, whatsappChatId string, participants []waha.GroupParticipant) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if participants == nil {
		stmt := `UPDATE groups 
				 SET title = $1, whatsapp_chat_id = $2 
				 WHERE id = $3`

		cmd, err := m.DB.Exec(ctx, stmt, title, whatsappChatId, id)
		if err != nil {
			return err
		}

		if cmd.RowsAffected() == 0 {
			return NotFoundError
		}

		return nil
	}

	tx, err := m.DB.Begin(ctx)

	defer tx.Rollback(ctx)

	if err != nil {
		return err
	}

	stmt := `UPDATE groups 
			 SET title = $1, whatsapp_chat_id = $2 
			 WHERE id = $3 
			 RETURNING id, title, whatsapp_chat_id, created_at`

	var group GroupDB

	err = tx.QueryRow(ctx, stmt, title, whatsappChatId, id).Scan(&group.ID, &group.Title, &group.WhatsappChatId, &group.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return NotFoundError
	}

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM group_members WHERE group_id = $1`, id)
	if err != nil {
		return err
	}

	for _, participant := range participants {
		err := m.AddParticipants(tx, ctx, group, participant)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *GroupModel) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM groups WHERE id = $1`

	cmd, err := m.DB.Exec(ctx, stmt, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return NotFoundError
	}

	return nil
}
