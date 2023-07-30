FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /app/bin/video-server /app/cmd/gateway/main.go

CMD ["/app/bin/video-server"]
