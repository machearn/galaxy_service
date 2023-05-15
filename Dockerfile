FROM golang:1.20-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .

EXPOSE 50051
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]