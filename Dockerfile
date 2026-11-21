FROM golang:1.27.4-alpine3.24 as builder
ARG VERSION="development"

WORKDIR /run

COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/cmd.Version=${VERSION}" -o /run ./cmd/remindme

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /run/remindme /run/remindme

USER 65532:65532
WORKDIR /run

CMD ["/run/remindme"]