# 1. build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 의존성 캐시
COPY go.mod go.sum ./
RUN go mod download

# 소스 복사
COPY . .

# static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o torchi

# 2. runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/torchi .


USER nonroot:nonroot
ENTRYPOINT ["/app/torchi"]
