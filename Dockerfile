FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /app

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o . ./...

FROM alpine:3.19

WORKDIR /usr/local/bin

COPY --from=builder /app/server .
COPY --from=builder /app/migrate .

COPY app.env .
COPY scripts/start.sh .
COPY sql/migrations/ ./sql/migrations/

CMD ["server"]

ENTRYPOINT [ "start.sh" ]

EXPOSE 8080