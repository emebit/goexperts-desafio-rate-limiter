FROM golang:1.22.4 as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./cmd/ -o rate-limiter

FROM scratch
WORKDIR /app
COPY --from=builder /app/cmd/.env /app/cmd/rate-limiter ./
ENTRYPOINT ["./rate-limiter"]