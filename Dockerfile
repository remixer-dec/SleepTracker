# Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Build backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o sleeptracker ./cmd/main.go

# Final distroless image
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=backend-builder /app/sleeptracker /app/sleeptracker
COPY --from=frontend-builder /app/frontend/dist /app/static

EXPOSE 8080

ENTRYPOINT ["/app/sleeptracker"]
CMD ["-port", "8080", "-db", "/app/data/sleeptracker.db", "-static", "/app/static"]
