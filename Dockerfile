# syntax=docker/dockerfile:1

FROM golang:1.25.4-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/rdvisit ./cmd/main.go

FROM alpine:3.22 AS app
RUN addgroup -S rdvisit && adduser -S rdvisit -G rdvisit
WORKDIR /app
COPY --from=builder /out/rdvisit /app/rdvisit
USER rdvisit
EXPOSE 3000
CMD ["/app/rdvisit"]

FROM migrate/migrate:v4.18.3 AS migrate
COPY db/migrations /migrations
ENTRYPOINT ["migrate"]
