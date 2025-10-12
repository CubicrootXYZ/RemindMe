FROM golang:1.25.2-alpine3.21@sha256:0ae17b3ad9583fcc9c2b195d12f2aa5dd1c18380d3827bd1a81c6e52aded353c as builder
ARG VERSION="development"

WORKDIR /run

COPY ./ ./
RUN go mod download
RUN go build -ldflags="-w -s -X github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/cmd.Version=${VERSION}" -o /run ./cmd/remindme

FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412
RUN apk update && apk upgrade
COPY --from=builder /run/remindme /run/
WORKDIR /run

CMD ["/run/remindme"]