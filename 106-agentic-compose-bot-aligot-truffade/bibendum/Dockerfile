FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN <<EOF
go mod tidy 
go build -o bibendum
EOF

FROM alpine:3.22.1
WORKDIR /app
COPY --from=builder /app/bibendum .

#CMD ["./bibendum"]