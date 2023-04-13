# Base go image
FROM golang:1.19-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN GOOS=linux CGO_ENABLED=0 go build -o frontendApp ./cmd/web

RUN chmod +x /app/frontendApp

#build tiny app
FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/static /app/static
COPY --from=builder /app/frontendApp /app   

CMD [ "/app/frontendApp" ]