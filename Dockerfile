FROM cubicrootxyz/matrix-go:latest

WORKDIR /run

COPY ./ ./
RUN go mod download
RUN go build -ldflags="-w -s" -o /run ./cmd/remindme

CMD ["/run/remindme"]
