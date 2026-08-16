FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /api ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 10001 appuser

WORKDIR /app
COPY --from=build /api /app/api

USER appuser

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/api"]
