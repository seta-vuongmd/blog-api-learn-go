# Build
FROM golang:1.23-alpine AS builder
WORKDIR /src
ENV CGO_ENABLED=0 GOOS=linux
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /out/blog-api-learn-go ./main.go

# Run (siêu nhỏ)
FROM scratch
COPY --from=builder /out/blog-api-learn-go /blog-api-learn-go
# nếu app gọi HTTPS, cần chứng chỉ:
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/blog-api-learn-go"]
