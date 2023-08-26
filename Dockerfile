FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app ./cmd/main.go && mv app ./cmd/

WORKDIR /app/cmd

EXPOSE 8080

CMD [ "./app" ]