# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN apk add wget
RUN wget https://github.com/eficode/wait-for/releases/download/v2.2.4/wait-for
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main . 
COPY --from=builder /app/wait-for .
COPY app.env .
COPY /src/db/migrations ./src/db/migrations

EXPOSE 8080
CMD [ "/app/main" ]
RUN chmod +x wait-for