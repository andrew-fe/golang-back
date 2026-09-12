# --- Etapa de compilación ---
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api

# --- Etapa final (imagen mínima) ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 appuser

COPY --from=builder /out/api /usr/local/bin/api

USER appuser
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/api"]
