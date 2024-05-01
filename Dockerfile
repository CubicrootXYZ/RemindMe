FROM cubicrootxyz/matrix-go:1.22 as builder
ARG VERSION="development"

WORKDIR /run

COPY ./ ./
RUN go mod download
RUN go build -ldflags="-w -s -X github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/cmd.Version=${VERSION}}" -o /run ./cmd/remindme

FROM ubuntu:22.04 
RUN apt update && apt upgrade -y && \
    apt install -y gcc && \
    apt install libolm-dev npm -y
COPY --from=builder /run/remindme /run/
WORKDIR /run

CMD ["/run/remindme"]