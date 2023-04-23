# polykube

## PostgreSQL

### Reference

* [**[Auth.js]** / API Reference / Database Adapters](https://authjs.dev/reference/adapters)
* [**[NextAuth.js (v3)]** / Documentation / Database Adapters / TypeORM / Postgres](https://next-auth.js.org/v3/adapters/typeorm/postgres)

### Run local development environment

```bash
$ cd polykube/.local
$ docker-compose up -d
```

### Create Databases

```postgresql
CREATE TABLE IF NOT EXISTS "nextauth_users" (
  "id" UUID,
  "name" VARCHAR(255),
  "email" VARCHAR(255),
  "email_verified" TIMESTAMP WITH TIME ZONE,
  "image" VARCHAR(255),
  UNIQUE ("email"),
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "nextauth_accounts" (
  "id" UUID,
  "type" VARCHAR(255) NOT NULL,
  "provider" VARCHAR(255) NOT NULL,
  "provider_account_id" VARCHAR(255) NOT NULL,
  "refresh_token" TEXT,
  "access_token" TEXT,
  "expires_at" INTEGER,
  "token_type" VARCHAR(255),
  "scope" VARCHAR(255),
  "id_token" TEXT,
  "session_state" VARCHAR(255),
  "user_id" UUID REFERENCES "nextauth_users" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
  PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS nextauth_accounts_provider_index ON nextauth_accounts (provider, provider_account_id);

CREATE TABLE IF NOT EXISTS "nextauth_sessions" (
  "id" UUID,
  "expires" TIMESTAMP WITH TIME ZONE NOT NULL,
  "session_token" VARCHAR(255) NOT NULL,
  "user_id" UUID REFERENCES "nextauth_users" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE ("session_token"),
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "nextauth_verification_token" (
  "identifier" VARCHAR(255) NOT NULL,
  "expires" TIMESTAMP WITH TIME ZONE NOT NULL,
  "token" VARCHAR(255) NOT NULL ,
  PRIMARY KEY ("token")
);
```

---

## Libraries

### Client

* Next.js
* NextAuth.js
* React Query

### Server

* Gin-Gonic
* Logrus
* Swag
