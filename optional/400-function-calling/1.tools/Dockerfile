FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY main.go .
COPY go.mod .

RUN <<EOF
go mod tidy 
go build -o function-calling
EOF

FROM scratch
WORKDIR /app
COPY --from=builder /app/function-calling .

CMD ["./function-calling"]