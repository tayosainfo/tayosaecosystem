FROM golang:1.22-alpine AS builder
WORKDIR /app

# Build api-gateway
COPY services/api-gateway-service/ ./api-gateway/
RUN cd api-gateway && go mod download && go build -ldflags="-s -w" -o /app/gateway .

# Build user-service
COPY services/user-service/ ./user-service/
RUN cd user-service && go mod download && go build -ldflags="-s -w" -o /app/user-server .

# Build object-storage-service
COPY services/object-storage-service/ ./object-storage-service/
RUN cd object-storage-service && go mod download && go build -ldflags="-s -w" -o /app/object-server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/gateway .
COPY --from=builder /app/user-server .
COPY --from=builder /app/object-server .

# Start script: runs user-service on 8081, object-storage on 8015, gateway on Render's $PORT
COPY <<'EOF' /app/start.sh
#!/bin/sh
# Render injects $PORT dynamically — use it for the gateway.
# Fall back to 8080 for local docker runs.
GATEWAY_PORT="${PORT:-8080}"

PORT=8081 ./user-server &
PORT=8015 ./object-server &
sleep 2
PORT=${GATEWAY_PORT} \
  USER_SERVICE_URL=http://localhost:8081 \
  OBJECT_STORAGE_SERVICE_URL=http://localhost:8015 \
  exec ./gateway
EOF

RUN chmod +x /app/start.sh
EXPOSE 10000
CMD ["/app/start.sh"]
