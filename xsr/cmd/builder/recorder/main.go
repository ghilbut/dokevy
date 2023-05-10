package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	// external
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v52/github"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func main() {

	// github
	webhookSecretKey := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if webhookSecretKey == "" {
		log.Fatal("github webhook secret does not exist")
	}

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

	// gin-gonic
	r := gin.New()
	r.Use(
		gin.LoggerWithWriter(
			log.StandardLogger().Out,
			"/healthz",
		),
		gin.Recovery(),
		func(ctx *gin.Context) {
			ctx.Set("DB", db)
			ctx.Set("WebhookSecretKey", webhookSecretKey)
		},
	)
	r.GET("/healthz", health)
	r.POST("/github/webhooks", handle)
	r.Run()
}

func health(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func handle(ctx *gin.Context) {
	db := ctx.MustGet("DB").(*gorm.DB)
	key := ctx.MustGet("WebhookSecretKey").(string)

	payload, err := github.ValidatePayload(ctx.Request, []byte(key))
	if err != nil {
		log.Fatal(err)
	}

	eventType := github.WebHookType(ctx.Request)
	event, err := github.ParseWebHook(eventType, payload)
	if err != nil {
		log.Fatal(err)
	}

	switch event := event.(type) {
	// Branch or tag creation
	case *github.CreateEvent:
		e := (*github.CreateEvent)(event)
		log.Tracef("%-18s : %s \n\t%s \n\t%s", "CreateEvent", e.GetRepo().GetFullName(), e.GetRefType(), e.GetRef())
		save(db, eventType, e.Repo.GetFullName(), string(payload))
		process(event)

	// Branch or tag deletion
	case *github.DeleteEvent:
		e := (*github.DeleteEvent)(event)
		log.Tracef("%-18s : %s \n\t%s \n\t%s", "DeleteEvent", e.GetRepo().GetFullName(), e.GetRefType(), e.GetRef())
		save(db, eventType, e.Repo.GetFullName(), string(payload))
		process(event)

	// Pull requests
	case *github.PullRequestEvent:
		e := (*github.PullRequestEvent)(event)
		log.Tracef("%-18s : %s \n\t[%d] %s", "PullRequestEvent", e.GetRepo().GetFullName(), e.GetNumber(), e.GetAction())
		save(db, eventType, e.Repo.GetFullName(), string(payload))
		process(event)

	// Pushes
	case *github.PushEvent:
		e := (*github.PushEvent)(event)
		log.Tracef("%-18s : %s \n\t%s \n\t%s", "Push", e.GetRepo().GetFullName(), e.GetRef(), e.GetAction())
		save(db, eventType, e.Repo.GetFullName(), string(payload))
		process(event)

	default:
		log.Infof("Ignore '%s' event type", eventType)
	}

	ctx.Status(http.StatusOK)
}

func process(event interface{}) {

}

func save(db *gorm.DB, repo, eventType, payload string) {
	tx := db.Select("Repository", "EventType", "Payload", "RequestedAt", "Status").Create(&entity{
		Repository:  repo,
		EventType:   eventType,
		Payload:     payload,
		RequestedAt: time.Now(),
		Status:      StatusRequested,
	})
	if tx.Error != nil {
		log.Fatal(tx.Error)
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
