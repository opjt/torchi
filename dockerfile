# 1. build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o torchi
    

# 2. runtime stage 
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/torchi .


USER nonroot:nonroot
ENTRYPOINT ["/app/torchi"]
