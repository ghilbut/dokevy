import { initTRPC } from '@trpc/server';
import superjson from "superjson";

const t = initTRPC.create({
    transformer: superjson,
});

export const middleware = t.middleware;
export const procedure = t.procedure;
export const router = t.router;
export const routers = t.mergeRouters;
