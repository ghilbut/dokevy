'use strict';

import {Pool, PoolConfig, QueryResult, QueryResultRow} from 'pg'
import {Awaitable} from 'next-auth';
import {AdapterAccount, AdapterSession, AdapterUser, DefaultAdapter, VerificationToken} from 'next-auth/adapters';
import {list} from "postcss";

const AccountTableName: string           = 'nextauth_accounts';
const SessionTableName: string           = 'nextauth_sessions';
const UserTableName: string              = 'nextauth_users';
const VerificationTokenTableName: string = 'nextauth_verification_tokens';

export default function NextAuth(): DefaultAdapter {

    // https://node-postgres.com/features/connecting#environment-variables
    // https://node-postgres.com/apis/client
    // https://node-postgres.com/apis/pool
    const pool: Pool = new Pool({
        // client
        application_name: 'dokevy',
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
    });

    function doQuery<T extends QueryResultRow>(query: { text: string, values: Array<any> }): Promise<T> {
        return new Promise<T>(async (resolve, reject) => {
            try {
                const res: QueryResult<T> = await pool.query(query);
                resolve(res.rows[0]);
            } catch (err) {
                reject(err);
            }
        });
    }

    function doQueryWithEmpty<T extends QueryResultRow, E extends null|undefined>(query: any, empty: E): Promise<T|E> {
        return new Promise<T | E>(async (resolve, reject) => {
            try {
                const res: QueryResult<T> = await pool.query(query);
                resolve(res.rowCount == 1 ? res.rows[0] : empty);
            } catch (err) {
                reject(err);
            }
        });
    }

    return {
        createUser(user: Omit<AdapterUser, 'id'>): Awaitable<AdapterUser> {
            const query = {
                text: `INSERT INTO ${UserTableName} (name, email, email_verified, image) VALUES ($1, $2, $3, $4) RETURNING *;`,
                values: [user.name, user.email, user.emailVerified || new Date(), user.image || null],
            }
            return doQuery<AdapterUser>(query);
        },

        getUser(id: string): Awaitable<AdapterUser | null> {
            const query = {
                text: `SELECT * FROM ${UserTableName} WHERE id = $1;`,
                values: [id],
            }
            return doQueryWithEmpty<AdapterUser, null>(query, null);
        },

        getUserByEmail(email: string): Awaitable<AdapterUser | null> {
            const query = {
                text: `SELECT * FROM ${UserTableName} WHERE email = $1;`,
                values: [email],
            }
            return doQueryWithEmpty<AdapterUser, null>(query, null);
        },

        getUserByAccount(
            providerAccountId: Pick<AdapterAccount, 'provider' | 'providerAccountId'>
        ): Awaitable<AdapterUser | null> {
            const query = {
                text: `SELECT * FROM ${UserTableName} WHERE id = 
                       (SELECT user_id FROM ${AccountTableName} WHERE provider = $1 AND provider_account_id = $2);`,
                values: [providerAccountId.provider, providerAccountId.providerAccountId],
            }
            return doQueryWithEmpty<AdapterUser, null>(query, null);
        },

        updateUser(user: Partial<AdapterUser> & Pick<AdapterUser, 'id'>): Awaitable<AdapterUser> {
            const query = {
                text: `UPDATE ${UserTableName} SET name = $2, email = $3, email_verified = $4, image = $5 WHERE id = $1 RETURNING *;`,
                values: [user.id, user.name, user.email, user.emailVerified, user.image],
            }
            return doQuery<AdapterUser>(query);
        },

        deleteUser(userId: string): Promise<void> | Awaitable<AdapterUser | null | undefined> {
            const query = {
                text: `DELETE FROM ${UserTableName} WHERE id = $1 RETURNING *;`,
                values: [userId],
            }
            return doQueryWithEmpty<AdapterUser, null>(query, null);
        },

        linkAccount(account: AdapterAccount): Promise<void> | Awaitable<AdapterAccount | null | undefined> {
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
            return doQueryWithEmpty<AdapterAccount, null>(query, null);
        },

        unlinkAccount(
            providerAccountId: Pick<AdapterAccount, 'provider' | 'providerAccountId'>
        ): Promise<void> | Awaitable<AdapterAccount | undefined> {
            const query = {
                text: `DELETE FROM ${AccountTableName} WHERE provider = $1 AND provider_account_id = $2 RETURNING *;`,
                values: [providerAccountId.provider, providerAccountId.providerAccountId],
            }
            return doQueryWithEmpty<AdapterAccount, undefined>(query, undefined);
        },

        createSession(session: { sessionToken: string, userId: string, expires: Date }): Awaitable<AdapterSession> {
            const query = {
                text: `INSERT INTO ${SessionTableName} (expires, session_token, user_id) VALUES ($1, $2, $3) RETURNING *;`,
                values: [session.expires, session.sessionToken, session.userId],
            }

            return new Promise<AdapterSession>(async (resolve, reject) => {
                try {
                    interface Row {
                        session_token: string;
                        user_id: string;
                        expires: Date;
                    }
                    const res: QueryResult<Row> = await pool.query(query);
                    const row = res.rows[0];
                    resolve({
                        sessionToken: row.session_token,
                        userId:       row.user_id,
                        expires:      row.expires,
                    });
                } catch (err) {
                    reject(err);
                }
            });
        },

        getSessionAndUser(sessionToken: string): Awaitable<{ session: AdapterSession, user: AdapterUser } | null> {
            const query = {
                text: `SELECT * FROM ${SessionTableName} s FULL OUTER JOIN ${UserTableName} u ON s.user_id = u.id WHERE s.session_token = $1;`,
                values: [sessionToken],
            }

            return new Promise<{ session: AdapterSession, user: AdapterUser } | null>(async (resolve, reject) => {
                try {
                    const res: QueryResult<any> = await pool.query(query);
                    const row = res.rows[0];
                    if (res.rowCount == 1) {
                        resolve({
                            session: {
                                expires:      row['expires'],
                                sessionToken: row['session_token'],
                                userId:       row['user_id'],
                            },
                            user: {
                                id:            row['id'],
                                name:          row['name'],
                                email:         row['email'],
                                emailVerified: row['email_verified'],
                                image:         row['image'],
                            }
                        });
                    } else {
                        resolve(null);
                    }
                } catch (err) {
                    reject(err);
                }
            });
        },

        updateSession(
            session: Partial<AdapterSession> & Pick<AdapterSession, 'sessionToken'>
        ): Awaitable<AdapterSession | null | undefined> {
            const query = {
                text: `UPDATE ${SessionTableName} SET expires = $1, user_id = $2 WHERE session_token = $3 RETURNING *;`,
                values: [session.expires, session.userId, session.sessionToken],
            }
            return doQueryWithEmpty<AdapterSession, null>(query, null);
        },

        deleteSession(sessionToken: string): Promise<void> | Awaitable<AdapterSession | null | undefined> {
            const query = {
                text: `DELETE FROM ${SessionTableName} WHERE session_token = $1 RETURNING *;`,
                values: [sessionToken],
            }
            return doQueryWithEmpty<AdapterSession, null>(query, null);
        },

        createVerificationToken(verificationToken: VerificationToken): Awaitable<VerificationToken | null | undefined> {
            const query = {
                text: `INSERT INTO ${VerificationTokenTableName} (identifier, expires, token) VALUES ($1, $2, $3) RETURNING *;`,
                values: [verificationToken.identifier, verificationToken.expires, verificationToken.token],
            }
            return doQueryWithEmpty<VerificationToken, null>(query, null);
        },

        useVerificationToken(params: { identifier: string, token: string }): Awaitable<VerificationToken | null> {
            const query = {
                text: `DELETE FROM ${VerificationTokenTableName} WHERE identifier = $1 AND token = $2 RETURNING *;`,
                values: [params.identifier, params.token],
            }
            return doQueryWithEmpty<VerificationToken, null>(query, null);
        },
    };
};
