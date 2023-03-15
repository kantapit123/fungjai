FROM golang:1.20 AS build

WORKDIR /app

COPY . .
RUN go mod download

ARG CHANNEL_SECRET
ARG CHANNEL_ACCESS_TOKEN
ARG BOOTSTRAP_SERVERS

RUN go build -o fungjai .

## Deploy
FROM alpine:latest 

WORKDIR /app

COPY --from=build /app/fungjai .

EXPOSE 1323

CMD ["sleep", "infinity"]
