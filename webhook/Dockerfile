FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN go mod download

ARG CHANNEL_SECRET
ARG CHANNEL_ACCESS_TOKEN
ARG BOOTSTRAP_SERVERS

# RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux

# # RUN go build -netgo fungjai .

# RUN go build -tags netgo -a -v fungjai .

RUN CGO_ENABLED=1 go build -tags netgo -a -v fungjai .

## Deploy
FROM alpine:3.15 

WORKDIR /app

COPY --from=build /app/fungjai .

EXPOSE 1323

# CMD [ "./fungjai" ]

ENTRYPOINT [ "./fungjai" ]
