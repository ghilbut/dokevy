package main

import (
	"fmt"
	"os"
	"time"

	// external
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func main() {

	// gorm
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("PG_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("PG_PASSWORD")
	database := os.Getenv("PG_DATABASE")
	if database == "" {
		database = "postgres"
	}
	format := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Seoul"
	dsn := fmt.Sprintf(format, host, port, user, password, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var rows []*entity
	tx := db.Find(&rows)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	for _, row := range rows {
		log.Infof("[%s] (%s) %v", row.Repository, row.EventType, row.RequestedAt)
	}

}

type Status string

const (
	StatusRequested  Status = "requested"
	StatusProcessing Status = "processing"
	StatusSucceed    Status = "succeed"
	StatusFailed     Status = "failed"
)

type entity struct {
	ID          int64      `gorm:"column:id"`
	Repository  string     `gorm:"column:repository"`
	EventType   string     `gorm:"column:type"`
	Payload     string     `gorm:"column:payload"`
	RequestedAt time.Time  `gorm:"column:requested"`
	ProcessedAt *time.Time `gorm:"column:processed"`
	CompletedAt *time.Time `gorm:"column:completed"`
	Status      Status     `gorm:"column:status"`
}

func (entity) TableName() string {
	return "github_webhook_events"
}
