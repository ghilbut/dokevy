ARG BUILDPLATFORM="linux/amd64"

FROM --platform=$BUILDPLATFORM golang:1.20.4 as builder

RUN mkdir /app
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o builder ./cmd/builder


FROM --platform=$BUILDPLATFORM scratch as release
LABEL co.ultary.image.authors="ghilbut@gmail.com"

ENV GOMAXPROCS=1
EXPOSE 8080
COPY --from=builder /app/builder /usr/local/bin/

CMD ["builder"]
