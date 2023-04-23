// NOTE(ghilbut): https://next-auth.js.org/v3/adapters/typeorm/postgres
package models

import "time"

type Account struct {
	ID                   uint64    `gorm:"column:id,type:BIGSERIAL,primaryKey"`
	CompoundID           string    `gorm:"column:compound_id,type:VARCHAR(255) NOT NULL"`
	UserID               int32     `gorm:"column:user_id,type:INTEGER NOT NULL"`
	provider_type        string    `gorm:"column:provider_type,type:VARCHAR(255) NOT NULL"`
	provider_id          string    `gorm:"column:provider_id,type:VARCHAR(255) NOT NULL"`
	provider_account_id  string    `gorm:"column:provider_account_id,type:VARCHAR(255) NOT NULL"`
	refresh_token        string    `gorm:"column:refresh_token,type:TEXT"`
	access_token         string    `gorm:"column:access_token,type:TEXT"`
	access_token_expires time.Time `gorm:"column:access_token_expires,type:TIMESTAMPTZ"`
	created_at           time.Time `gorm:"column:created_at,type:TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	updated_at           time.Time `gorm:"column:updated_at,type:TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (Account) TableName() string {
	return "nextauth_accounts"
}

type Session struct {
	id            uint64    `gorm:"column:id,type:BIGSERIAL,primaryKey"`
	user_id       int32     `gorm:"column:user_id,type:INTEGER NOT NULL"`
	expires       time.Time `gorm:"column:expires,type:TIMESTAMPTZ NOT NULL"`
	session_token string    `gorm:"column:session_token,type:VARCHAR(255) NOT NULL"`
	access_token  string    `gorm:"column:access_token,type:VARCHAR(255) NOT NULL"`
	created_at    time.Time `gorm:"column:created_at,type:TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	updated_at    time.Time `gorm:"column:updated_at,type:TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (Session) TableName() string {
	return "nextauth_sessions"
}
