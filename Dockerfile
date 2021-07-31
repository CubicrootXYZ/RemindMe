FROM golang:1.16.3-alpine

WORKDIR /app

COPY * ./
RUN go mod download
RUN go build -o /run ./cmd/remindme

CMD ["/run"]
