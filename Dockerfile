FROM golang:1.20-alpine

EXPOSE 2620

RUN mkdir -pv /usr/src/app
COPY . /usr/src/app
WORKDIR /usr/src/app

RUN go build

ENTRYPOINT ["./gin-examples"]
