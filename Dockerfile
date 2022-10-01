FROM golang:alpine as builder

LABEL maintainer="Riszky MF <riszkymfahreza@gmail.com>"

RUN apk update && apk add --no-cache git
RUN apk add build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . ./

RUN cd cmd && go build -o ../server && cd ..
RUN ["chmod","+x","./server"]

CMD ["./server"]