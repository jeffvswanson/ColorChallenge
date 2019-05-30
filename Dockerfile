FROM golang:1.12 AS builder

RUN mkdir /app
WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest AS certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=builder /app/main .
COPY --from=builder /app/input.txt .
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["./main"]