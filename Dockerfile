# syntax=docker/dockerfile:1.5.2

ARG BUILDPLATFORM="linux/amd64"
ARG NODETAG="18.16-slim"

FROM --platform=$BUILDPLATFORM node:$NODETAG AS builder

COPY fastify  /app/fastify
WORKDIR /app/fastify
RUN yarn install \
 && yarn build

COPY react.js /app/react.js
WORKDIR /app/react.js
RUN yarn install \
 && yarn build


FROM --platform=$BUILDPLATFORM node:$NODETAG AS runtime

COPY fastify /app/fastify
WORKDIR /app/fastify
RUN yarn install --production


FROM --platform=$BUILDPLATFORM node:$NODETAG AS client
LABEL com.ghilbut.image.authors="ghilbut@gmail.com"

ENV NODE_ENV=production \
    REACT_MODE=static \
    REACT_ROOT=/app/htdocs
EXPOSE 3030
WORKDIR /app/fastify

COPY --from=builder /app/react.js/build       /app/htdocs
COPY --from=builder /app/fastify/bin          /app/fastify
COPY --from=runtime /app/fastify/node_modules /app/fastify/node_modules

CMD [ "/bin/sh", "-c", "node /app/fastify/index.js" ]
