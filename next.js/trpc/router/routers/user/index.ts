import { z } from 'zod';
import { router, procedure } from '../../trpc';

export const userRouter = router({
    list: procedure.query(() => {
        // [..]
        return [];
    }),
});
