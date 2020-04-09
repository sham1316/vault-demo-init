FROM golang:1.13.4 AS builder
RUN go version

COPY . /app/
WORKDIR /app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM alpine:3.11.5
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .

EXPOSE 80
ENTRYPOINT ["./app"]