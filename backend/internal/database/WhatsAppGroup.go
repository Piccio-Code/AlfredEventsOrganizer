package database

import "github.com/jackc/pgx/v5/pgxpool"

type WhatsAppGroupModel struct {
	DB *pgxpool.Pool
}

type WhatsAppGroupCreate struct {
}

func (m *WhatsAppGroupModel) Create() {

}
