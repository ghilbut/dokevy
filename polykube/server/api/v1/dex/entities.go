package dex

type Action string

const (
	Create Action = "create"
	Update Action = "action"
	Delete Action = "delete"
)

type AuditEntity struct {
	Author   string `gorm:"column:author"`
	Action   Action `gorm:"column:action"`
	ClientID string `gorm:"column:client_id"`
	Message  string `gorm:"column:message"`
}
