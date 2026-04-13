FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gophershield ./cmd/main.go

FROM gcr.io/distroless/static
COPY --from=builder /app/gophershield /gophershield
EXPOSE 8080 9090
ENTRYPOINT ["/gophershield"]