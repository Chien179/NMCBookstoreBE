# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN apk add curl
RUN curl -L https://github.com/swaggo/swag/releases/download/v1.8.12/swag_1.8.12_Linux_x86_64.tar.gz | tar xvz
RUN ./swag init
RUN apk add wget
RUN wget https://github.com/eficode/wait-for/releases/download/v2.2.4/wait-for
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main . 
COPY --from=builder /app/wait-for .
COPY app.env .
COPY db/migrations ./db/migrations

EXPOSE 8080
CMD [ "/app/main" ]
RUN chmod +x wait-for