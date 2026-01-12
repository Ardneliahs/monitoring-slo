# -------- Build stage --------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod dservice.go ./
RUN go mod download

# Build static binary disables C static links builds for linux 64 bit amd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o dservice

# -------- Runtime stage --------
#distroless -- does not ship with system binaries and libraries. Less attack surface.
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/dservice .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/dservice"]
