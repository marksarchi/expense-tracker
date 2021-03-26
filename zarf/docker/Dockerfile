FROM golang:1.15.6 as builder

ENV CGO_ENABLED 0

RUN mkdir -p /expense-tracker

COPY go.* /expense-tracker/
WORKDIR /expense-tracker


RUN go mod download

WORKDIR /expense-tracker

COPY . .

WORKDIR /expense-tracker/app

RUN go build

#Run go binary in alpine

FROM alpine:latest AS production

COPY  --from=builder /expense-tracker/app/app /expense-tracker/app
WORKDIR /expense-tracker
CMD ["./app"]

#Notice that we’ve added the go.* files and download the modules before adding the rest of the source.
# This allows Docker to cache the modules as it will only rerun these steps if the go.* files change.

#ARG instruction lets Dockerfile authors define values that users can set at build-time using the --build-arg flag:





