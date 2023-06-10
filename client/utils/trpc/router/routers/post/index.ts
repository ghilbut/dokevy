import { z } from 'zod';
import { router, procedure } from '../../trpc';

export const postRouter = router({
    create: procedure
        .input(
            z.object({
                title: z.string(),
            }),
        )
        .mutation((opts) => {

        }),
    list: procedure.query(() => {
        // ...
        return [];
    }),
});
