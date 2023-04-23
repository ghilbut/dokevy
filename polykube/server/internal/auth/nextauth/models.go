package nextauth

import "time"

const (
	AccountTableName           = "nextauth_accounts"
	SessionTableName           = "nextauth_sessions"
	UserTableName              = "nextauth_users"
	VerificationTokenTableName = "nextauth_verification_tokens"
)

type AccountEntity struct {
	ID                string `gorm:"column:id"`
	UserID            string `gorm:"column:user_id"`
	Type              string `gorm:"column:type"`
	Provider          string `gorm:"column:provider"`
	ProviderAccountId string `gorm:"column:provider_account_id"`
	RefreshToken      string `gorm:"column:refresh_token"`
	AccessToken       string `gorm:"column:access_token"`
	ExpiresAt         int    `gorm:"column:expires_at"`
	TokenType         string `gorm:"column:token_type"`
	Scope             string `gorm:"column:scope"`
	IDToken           string `gorm:"column:id_token"`
	SessionState      string `gorm:"column:session_state"`
	OAuthTokenSecret  string `gorm:"column:oauth_token_secret"`
	OAuthToken        string `gorm:"column:oauth_token"`
}

func (AccountEntity) TableName() string {
	return AccountTableName
}

type SessionEntity struct {
	ID           string    `gorm:"column:id"`
	Expires      time.Time `gorm:"column:expires"`
	SessionToken string    `gorm:"column:session_token"`
	UserID       string    `gorm:"column:user_id"`
}

func (SessionEntity) TableName() string {
	return SessionTableName
}

type UserEntity struct {
	ID            string    `gorm:"column:id"`
	Name          string    `gorm:"column:name"`
	Email         string    `gorm:"column:email"`
	EmailVerified time.Time `gorm:"column:email_verified"`
	Image         string    `gorm:"column:image"`
}

func (UserEntity) TableName() string {
	return UserTableName
}

type VerificationTokenEntity struct {
	identifier string    `gorm:"column:identifier"`
	token      string    `gorm:"column:token"`
	expires    time.Time `gorm:"column:expires"`
}

func (VerificationTokenEntity) TableName() string {
	return VerificationTokenTableName
}
