
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /main ./cmd/node

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /main ./main

ENV PORT=8080

EXPOSE 8080

CMD ["./main"]
