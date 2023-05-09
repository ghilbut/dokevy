# XSR

```postgresql
CREATE DATABASE xsr OWNER postgres;
```

```postgresql
CREATE TABLE IF NOT EXISTS "github_webhook_events"
(
    "id"        BIGSERIAL,
    "payload"   TEXT,
    "requested" TIMESTAMP WITH TIME ZONE NOT NULL,
    "processed" TIMESTAMP WITH TIME ZONE,
    "completed" TIMESTAMP WITH TIME ZONE,
    "status"    VARCHAR(255),
    PRIMARY KEY ("id")
);
```
