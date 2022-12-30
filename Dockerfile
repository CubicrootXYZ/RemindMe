FROM cubicrootxyz/matrix-go:1.19 as builder

WORKDIR /run

COPY ./ ./
RUN go mod download
RUN go build -ldflags="-w -s" -o /run ./cmd/remindme

FROM ubuntu:22.04 
RUN apt update && apt upgrade -y && \
    apt install -y gcc && \
    apt install libolm-dev npm -y
COPY --from=builder /run/remindme /run/
WORKDIR /run

CMD ["/run/remindme"]