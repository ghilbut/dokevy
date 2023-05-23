# PostgreSQL

## NextAuth.js

* [**[Auth.js]** / API Reference / Database Adapters](https://authjs.dev/reference/adapters)
* [**[NextAuth.js (v3)]** / Documentation / Database Adapters / TypeORM / Postgres](https://next-auth.js.org/v3/adapters/typeorm/postgres)

**Create Databases**

```postgresql
CREATE TABLE IF NOT EXISTS "nextauth_users" (
  "id" BIGSERIAL,
  "name" VARCHAR(255),
  "email" VARCHAR(255),
  "email_verified" TIMESTAMP WITH TIME ZONE,
  "image" VARCHAR(255),
  UNIQUE ("email"),
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "nextauth_accounts" (
  "id" BIGSERIAL,
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
  "user_id" BIGSERIAL REFERENCES "nextauth_users" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
  PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS nextauth_accounts_provider_index ON nextauth_accounts (provider, provider_account_id);

CREATE TABLE IF NOT EXISTS "nextauth_sessions" (
  "id" BIGSERIAL,
  "expires" TIMESTAMP WITH TIME ZONE NOT NULL,
  "session_token" VARCHAR(255) NOT NULL,
  "user_id" BIGSERIAL REFERENCES "nextauth_users" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
  UNIQUE ("session_token"),
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "nextauth_verification_tokens" (
  "identifier" VARCHAR(255) NOT NULL,
  "expires" TIMESTAMP WITH TIME ZONE NOT NULL,
  "token" VARCHAR(255) NOT NULL ,
  PRIMARY KEY ("token")
);

CREATE INDEX IF NOT EXISTS nextauth_verification_token_index ON nextauth_verification_tokens (identifier, token);
```

## Teraform

**Create Databases**

```postgresql

CREATE TABLE IF NOT EXISTS "dokevy_secret_shared_secrets" (
  "id"         BIGSERIAL,
  "name"       VARCHAR(255) NOT NULL,
  "secrets"    JSON,
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  UNIQUE ("name"),
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "dokevy_secret_scoped_groups" (
  "id"         BIGSERIAL,
  "name"       VARCHAR(255) NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  UNIQUE ("name"),
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "dokevy_secret_scoped_secrets" (
  "id"         BIGSERIAL,
  "group"      BIGSERIAL REFERENCES "dokevy_secret_scoped_groups" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
  "secrets"    JSON,
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  UNIQUE ("group"),
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "dokevy_terraform_bots" (
  "id"          BIGSERIAL,
  "username"    VARCHAR(255) NOT NULL,
  "password"    VARCHAR(64) NOT NULL,
  "description" TEXT,
  "created_by"  VARCHAR(255) NOT NULL,
  "created_at"  TIMESTAMP WITH TIME ZONE NOT NULL,
  "updated_at"  TIMESTAMP WITH TIME ZONE NOT NULL,
  UNIQUE ("username"),
  PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS dokevy_terraform_bot_creator_index ON dokevy_terraform_bots (created_by);
```
