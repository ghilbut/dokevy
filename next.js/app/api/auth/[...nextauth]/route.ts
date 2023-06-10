import NextAuth from 'next-auth';
import type { NextAuthOptions } from 'next-auth';
import Adapter from '~/adapters/next-auth';

export const authOptions: NextAuthOptions = {
    providers: [
        {
            id:   'dex',
            name: 'Dex',
            type: 'oauth',
            clientId: 'nextauthjs',
            clientSecret: 'nextauthjs',
            wellKnown: 'http://localhost:5556/.well-known/openid-configuration',
            authorization: { params: { scope: 'openid email profile' } },
            idToken: true,
            checks: ['pkce', 'state'],
            profile(profile) {
                return {
                    id: profile.sub,
                    name: profile.name,
                    email: profile.email,
                    image: profile.picture,
                }
            },
        },
    ],
    session: {
        strategy: 'database',
        maxAge: 30 * 24 * 60 * 60, // 30 days
        updateAge: 24 * 60 * 60, // 24 hours
    },
    adapter: Adapter(),
    debug: true
};

const handler = NextAuth(authOptions);
export { handler as GET, handler as POST };
