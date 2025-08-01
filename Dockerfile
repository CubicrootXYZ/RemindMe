FROM golang:1.24.5-alpine3.21 as builder
ARG VERSION="development"

WORKDIR /run

COPY ./ ./
RUN go mod download
RUN go build -ldflags="-w -s -X github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/cmd.Version=${VERSION}" -o /run ./cmd/remindme

FROM alpine:3.22
COPY --from=builder /run/remindme /run/
WORKDIR /run

CMD ["/run/remindme"]