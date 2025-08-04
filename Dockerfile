FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /app
RUN apk add --no-cache tzdata

COPY --from=builder /app/cmd/main.go ./
COPY --from=builder /app/app ./
COPY --from=builder /app/.env ./

EXPOSE 4242

CMD ["./app"]