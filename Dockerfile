FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/habrabot

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /out/habrabot /app/
ENV BOLTDB_PATH /data/rss.dat
CMD ./habrabot
