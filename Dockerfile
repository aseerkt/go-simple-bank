FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o . ./...

FROM alpine:3.19

WORKDIR /usr/local/bin

COPY --from=builder /app/server .

COPY app.env .
COPY sql/migrations/ ./sql/migrations/

CMD ["server"]

EXPOSE 8080