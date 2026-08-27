# Stage 1: Build Frontend Web Assets (Svelte 5 + Vite + Tailwind)
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Pure-Go Static Server Binary
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.* ./
RUN go mod download
COPY backend/ ./
# Copy built static assets into backend/src/dist for go:embed
COPY --from=frontend-builder /app/backend/src/dist ./src/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o led-swarm-server ./src

# Stage 3: Minimal Alpine Production Container (<25MB)
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /app/backend/led-swarm-server ./
EXPOSE 8080
ENV SERVER_MODE=1
ENTRYPOINT ["./led-swarm-server", "--server", "--port=8080"]
