FROM golang:1.20 AS build

WORKDIR /app

COPY . .
RUN go mod download

ARG CHANNEL_SECRET
ARG CHANNEL_ACCESS_TOKEN
ARG BOOTSTRAP_SERVERS


RUN go build -o /fungjai

EXPOSE 1323

CMD [ "/fungjai" ]