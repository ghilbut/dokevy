import { initTRPC } from "@trpc/server";
import superjson from "superjson";
import { middleware, procedure, router, routers } from './trpc'

import { exampleRouter } from './routers/example';
import { postRouter } from './routers/post';
import { userRouter } from './routers/user';

export const appRouter = routers(
    exampleRouter,
    router({
       post: postRouter,
       user: userRouter,
    }),
);

export type AppRouter = typeof appRouter;
