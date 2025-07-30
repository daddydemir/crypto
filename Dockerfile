FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app ./cmd/main.go && mv app ./cmd/ && mv .env /app/cmd/

WORKDIR /app/cmd

EXPOSE 4242

CMD [ "./app" ]