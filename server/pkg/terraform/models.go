package terraform

import "time"

const (
	SecretBackendTableName = "polykube_terraaform_secret_backends"
	SecretValueTableName   = "polykube_terraaform_secret_values"
	BotTableName           = "polykube_terraform_bots"
)

type SecretBackendEntity struct {
	ID        uint64    `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type SecretValueEntity struct {
	ID        uint64    `gorm:"column:id"`
	Secret    string    `gorm:"column:secret"`
	Key       string    `gorm:"column:key"`
	Value     string    `gorm:"column:value"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type BotEntity struct {
	ID          uint64    `gorm:"column:id"`
	Username    string    `gorm:"column:username"`
	Password    string    `gorm:"column:password"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}
