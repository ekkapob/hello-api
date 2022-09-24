FROM golang:1.18-alpine as builder
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM alpine
RUN apk add --no-cache curl
RUN apk add --no-cache bash
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]

