FROM golang:1.22.0-alpine AS builder
WORKDIR /app
ENV GO111MODULE=on
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o aws-key-scanner ./cmd/awsKeyhunter.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/aws-key-scanner .
USER nonroot:nonroot
ENTRYPOINT ["/app/aws-key-scanner"]
