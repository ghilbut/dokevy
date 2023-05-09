package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	// external
	"github.com/google/go-github/v52/github"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.SetLevel(log.TraceLevel)

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

	// fasthttp
	log.Info("RUN XSR Recorder on 0.0.0.0:8080")
	if err := fasthttp.ListenAndServe("0.0.0.0:8080", handle(db, webhookSecretKey)); err != nil {
		log.Fatal(err)
	}
}

// handle returns a fasthttp request handler
func handle(db *gorm.DB, webhookSecretKey string) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := ctx.Path()
		if string(path) == "/healthz" {
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBody([]byte("OK"))
			return
		}

		if !ctx.IsPost() {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			return
		}

		defer func() {
			if r := recover(); r != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				ctx.SetContentType("text/html; charset=utf-8")
				ctx.SetBody([]byte("Internal Server Error"))
			}
		}()

		const forServer = false
		var r http.Request
		if err := fasthttpadaptor.ConvertRequest(ctx, &r, forServer); err != nil {
			log.Error(fasthttp.StatusBadRequest, err)
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		payload, err := github.ValidatePayload(&r, []byte(webhookSecretKey))
		if err != nil {
			log.Fatal(err)
		}

		eventType := github.WebHookType(&r)
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
	}
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
