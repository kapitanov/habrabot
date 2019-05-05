FROM golang:1.12 AS builder
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 go build -o /out/habrabot

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /out/habrabot /app/
CMD ./habrabot
