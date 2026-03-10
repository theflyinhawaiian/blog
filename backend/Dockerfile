FROM golang:1.25 AS base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS dev
RUN go install github.com/air-verse/air@latest
EXPOSE 8080
CMD ["air"]

FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -o blog ./main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/blog .
EXPOSE 8080
CMD ["./blog"]
