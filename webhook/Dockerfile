FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY *.go .
# COPY producer/ .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -tags netgo -a -v fungjai .
# RUN CGO_ENABLED=0 go build -tags netgo -a -v -o /app/fungjai .
# RUN go build -tags musl -o fungjai .

## Deploy
FROM ubuntu:22.04

WORKDIR /app

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates

COPY --from=build /app/fungjai .

EXPOSE 1323

# CMD [ "sleep", "infinity" ]

# ENTRYPOINT [ "/app/fungjai" ]
CMD [ "/app/fungjai" ]
