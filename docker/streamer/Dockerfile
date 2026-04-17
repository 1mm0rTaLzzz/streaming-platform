FROM golang:1.23-alpine AS builder
# Go 1.23 is latest stable — works on linux/amd64 VPS
WORKDIR /build
COPY streamctl/go.mod ./
RUN go mod download
COPY streamctl/main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o streamctl .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium \
    xvfb \
    ffmpeg \
    dumb-init \
    ca-certificates \
    fonts-liberation \
    fonts-noto \
    libnss3 \
    libglib2.0-0 \
    libx11-xcb1 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxfixes3 \
    libxi6 \
    libxrandr2 \
    libxss1 \
    libxtst6 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/streamctl .
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8090
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/entrypoint.sh"]
