# builder
FROM golang:1.20-alpine AS builder

WORKDIR /home/app

RUN apk add --update make 

COPY . .
RUN go build -o bin/merchant-experience cmd/app/main.go

# app
FROM alpine:latest as app

WORKDIR /root/

COPY --from=builder /home/app/bin/merchant-experience .
COPY migrations migrations

# RUN chown root:root merchant-experience

EXPOSE 8000

CMD ["./merchant-experience"]