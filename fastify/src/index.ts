import path from 'path';
import fastify from 'fastify';
import fastifyProxy from '@fastify/http-proxy';
import fastifyStatic from '@fastify/static';
import fastifyWebsocket from '@fastify/websocket';
import { fastifyTRPCPlugin } from '@trpc/server/adapters/fastify';
import { appRouter } from './router';

const server = fastify({
    logger: true,
    maxParamLength: 4096,
});

server.register(fastifyWebsocket);
server.register(fastifyTRPCPlugin, {
    prefix: '/trpc',
    trpcOptions: { router: appRouter },
    useWSS: true,
});

server.get('/healthz', async (request, reply) => {
    return 'OK'
});

const reactMode = process.env.REACT_MODE || 'proxy';
switch (reactMode) {
    case 'none': {
        console.log('do not serve react.js');
        break;
    }
    case 'proxy': {
        server.register(fastifyProxy, {
            http2: false,
            upstream: 'http://127.0.0.1:3000',
            websocket: true,
            wsUpstream: 'ws://127.0.0.1:3000',
        });

        console.log('serve react development server with proxy');
        break;
    }
    case 'static': {
        const rootdir = process.env.REACT_ROOT || path.join(__dirname, '../../react.js/build');
        server.register(fastifyStatic, {
            root: rootdir,
        });
        server.setNotFoundHandler((req, reply) => {
            reply.sendFile('index.html', rootdir);
            reply.status(200);
        });

        console.log('serve static react files');
        break;
    }
    default: {
        console.error(`${reactMode} is unsupported environment`);
        process.exit(1);
    }
}

server.listen({
    host: '0.0.0.0',
    port: 3030,
}).then(address => {
    console.log(`server listening on ${address}`);
}).catch(err => {
    console.error('Error starting server:', err);
    process.exit(1);
});
