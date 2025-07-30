# ---- Build Stage ----
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# ---- Run Stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main ./main
RUN chmod +x main
COPY static ./static
COPY views ./views
COPY sql ./sql
COPY data ./data
EXPOSE 8080
CMD ["./main"]