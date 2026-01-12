# -------- Build stage --------
FROM golang:1.25.5-alpine AS builder

WORKDIR /monitor

COPY go.mod go.sum monitor.go ./
RUN go mod download

# Build static binary disables C static links builds for linux 64 bit amd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o monitor

# -------- Runtime stage --------
#distroless -- does not ship with system binaries and libraries. Less attack surface.
FROM gcr.io/distroless/base-debian12

WORKDIR /monitor

COPY --from=builder /monitor/monitor .

EXPOSE 8081

USER nonroot:nonroot

ENTRYPOINT ["/monitor/monitor --config=/monitor/config.yaml"]
