FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux

RUN cd cmd/cure && go build -a -ldflags '-extldflags "-static"' -o /app/bin/cure .

FROM scratch
COPY --from=builder  /app/bin/cure /
