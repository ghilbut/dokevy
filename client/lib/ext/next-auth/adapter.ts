'use strict';

import {Pool, PoolConfig, QueryResult} from 'pg'
import {Awaitable} from 'next-auth';
import {AdapterAccount, AdapterSession, AdapterUser, DefaultAdapter, VerificationToken} from 'next-auth/adapters';

const AccountTableName: string           = 'nextauth_accounts';
const SessionTableName: string           = 'nextauth_sessions';
const UserTableName: string              = 'nextauth_users';
const VerificationTokenTableName: string = 'nextauth_verification_tokens';

export default function Adapter(): DefaultAdapter {

    // https://node-postgres.com/features/connecting#environment-variables
    // https://node-postgres.com/apis/client
    // https://node-postgres.com/apis/pool
    const config: PoolConfig = {
        // client
        application_name: 'polykube',
        host: process.env.PGHOST || 'localhost',
        port: Number(process.env.PGPORT || '5432'),
        user: process.env.PGUSER || 'postgres',
        password: process.env.PGPASSWORD || 'postgrespw',
        database: process.env.PGDATABASE || 'postgres',
        keepAlive: true,
        connectionTimeoutMillis:              2 * 1000, //  2 * 1000 milliseconds =  2 seconds
        idle_in_transaction_session_timeout:  2 * 1000, // 10 * 1000 milliseconds = 10 seconds
        query_timeout:                       10 * 1000, // 10 * 1000 milliseconds = 10 seconds
        statement_timeout:                   60 * 1000, // 60 * 1000 milliseconds =  1 minutes
        // pool
        max: Number(process.env.PGMAX || '20'),
        min: Number(process.env.PGMIN || '5'),
    };
    const pool: Pool = new Pool(config);

    return {
        async createUser(user: Omit<AdapterUser, 'id'>): Awaitable<AdapterUser> {
            const query = {
                text: `INSERT INTO ${UserTableName} (name, email, email_verified, image) VALUES ($1, $2, $3, $4) RETURNING *;`,
                values: [user.name, user.email, user.emailVerified || new Date(), user.image || null],
            }

            const res: QueryResult<AdapterUser> = await pool.query(query);
            return res.rows[0];
        },

        async getUser(id: string): Awaitable<AdapterUser | null> {
            const query = {
                text: `SELECT * FROM ${UserTableName} WHERE id = $1;`,
                values: [id],
            }

            const res: QueryResult<AdapterUser> = await pool.query(query);
            return res.rowCount == 1 ? res.rows[0] : null;
        },

        async getUserByEmail(email: string): Awaitable<AdapterUser | null> {
            const query = {
                text: `SELECT * FROM ${UserTableName} WHERE email = $1;`,
                values: [email],
            }

            const res: QueryResult<AdapterUser> = await pool.query(query);
            return res.rowCount == 1 ? res.rows[0] : null;
        },

        async getUserByAccount(
            providerAccountId: Pick<AdapterAccount, 'provider' | 'providerAccountId'>
        ): Awaitable<AdapterUser | null> {
            const query = {
                text: `SELECT * FROM ${UserTableName} WHERE id = 
                       (SELECT user_id FROM ${AccountTableName} WHERE provider = $1 AND provider_account_id = $2);`,
                values: [providerAccountId.provider, providerAccountId.providerAccountId],
            }

            const res: QueryResult<AdapterUser> = await pool.query(query);
            return res.rowCount == 1 ? res.rows[0] : null;
        },

        async updateUser(user: Partial<AdapterUser> & Pick<AdapterUser, 'id'>): Awaitable<AdapterUser> {
            const query = {
                text: `UPDATE ${UserTableName} SET name = $2, email = $3, email_verified = $4, image = $5 WHERE id = $1 RETURNING *;`,
                values: [user.id, user.name, user.email, user.emailVerified, user.image],
            }

            const res: QueryResult<AdapterUser> = await pool.query(query);
            return res.rows[0];
        },

        async deleteUser(userId: string): Promise<void> | Awaitable<AdapterUser | null | undefined> {
            const query = {
                text: `DELETE FROM ${UserTableName} WHERE id = $1 RETURNING *;`,
                values: [providerAccountId.provider, providerAccountId.providerAccountId],
            }

            const res: QueryResult<AdapterUser> = await pool.query(query);
            return res.rowCount == 1 ? res.rows[0] : null;
        },

        async linkAccount(account: AdapterAccount): Promise<void> | Awaitable<AdapterAccount | null | undefined> {
            const query = {
                text: `INSERT INTO ${AccountTableName} (
                         type,
                         provider,
                         provider_account_id,
                         refresh_token,
                         access_token,
                         expires_at,
                         token_type,
                         scope,
                         id_token,
                         session_state,
                         user_id
                       )
                       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;`,
                values: [
                    account.type,
                    account.provider,
                    account.providerAccountId,
                    account.refresh_token || null,
                    account.access_token,
                    account.expires_at,
                    account.token_type,
                    account.scope || null,
                    account.id_token,
                    account.session_state,
                    account.userId,
                ],
            }

            const res: QueryResult<AdapterAccount> = await pool.query(query);
            return res.rowCount == 1 ? res.rows[0] : null;
        },

        async unlinkAccount(
            providerAccountId: Pick<AdapterAccount, 'provider' | 'providerAccountId'>
        ): Promise<void> | Awaitable<AdapterAccount | undefined> {
            const query = {
                text: `DELETE FROM ${AccountTableName} WHERE provider = $1 AND provider_account_id = $2 RETURNING *;`,
                values: [providerAccountId.provider, providerAccountId.providerAccountId],
            }

            const res: QueryResult<AdapterAccount> = await pool.query(query);
            return res.rowCount == 1 ? res.rows[0] : null;
        },

        async createSession(session: {
            sessionToken: string;
            userId: string;
            expires: Date;
        }): Awaitable<AdapterSession> {
            const query = {
                text: `INSERT INTO ${SessionTableName} (expires, session_token, user_id) VALUES ($1, $2, $3) RETURNING *;`,
                values: [session.expires, session.sessionToken, session.userId],
            }

            const res: QueryResult<AdapterSession> = await pool.query(query);
            const row = res.rows[0];
            return {
                expires: new Date(row['expires']),
                sessionToken: row['session_token'],
                userId: row['user_id'],
            };
        },

        async getSessionAndUser(sessionToken: string): Awaitable<{
            session: AdapterSession;
            user: AdapterUser;
        } | null> {
            const query = {
                text: `SELECT * FROM ${SessionTableName} s FULL OUTER JOIN ${UserTableName} u ON s.user_id = u.id WHERE s.session_token = $1;`,
                values: [sessionToken],
            }

            const res: QueryResult<AdapterSession> = await pool.query(query);
            const row = res.rows[0];
            if (res.rowCount == 1) {
                return {
                    session: {
                        expires:       row['expires'],
                        session_token: row['session_token'],
                        user_id:       row['user_id'],
                    },
                    user: {
                        id:            row['id'],
                        name:          row['name'],
                        email:         row['email'],
                        emailVerified: row['email_verified'],
                        image:         row['image'],
                    }
                };
            }
            return null;
        },

        async updateSession(
            session: Partial<AdapterSession> & Pick<AdapterSession, 'sessionToken'>
        ): Awaitable<AdapterSession | null | undefined> {
            const query = {
                text: `UPDATE ${SessionTableName} SET expires = $1, user_id = $2 WHERE session_token = $3 RETURNING *;`,
                values: [session.expires, session.userId, session.sessionToken],
            }

            const res: QueryResult<AdapterSession> = await pool.query(query);
            return res.rowCount == 1 ? res.rows[0] : null;
        },

        async deleteSession(sessionToken: string): Promise<void> | Awaitable<AdapterSession | null | undefined> {
            const query = {
                text: `DELETE FROM ${SessionTableName} WHERE session_token = $1 RETURNING *;`,
                values: [sessionToken],
            }

            const res: QueryResult<AdapterSession> = await pool.query(query);
            return res.rowCount == 1 ? res.rows[0] : null;
        },

        async createVerificationToken(
            verificationToken: VerificationToken
        ): Awaitable<VerificationToken | null | undefined> {
            const query = {
                text: `INSERT INTO ${VerificationTokenTableName} (identifier, expires, token) VALUES ($1, $2, $3) RETURNING *;`,
                values: [verificationToken.identifier, verificationToken.expires, verificationToken.token],
            }

            const res:QueryResult<VerificationToken> = await pool.query(query);
            return res.rowCount == 1 ? res.rows[0] : null;
        },

        async useVerificationToken(params: {
            identifier: string;
            token: string;
        }): Awaitable<VerificationToken | null> {
            const query = {
                text: `DELETE FROM ${VerificationTokenTableName} WHERE identifier = $1 AND token = $2 RETURNING *;`,
                values: [identifier, token],
            }

            const res:QueryResult<VerificationToken> = await pool.query(query);
            return res.rowCount == 1 ? res.rows[0] : null;
        },
    };
};
