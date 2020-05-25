FROM golang:1.14.0-alpine

WORKDIR src/nats-docker-api

COPY . .

#RUN pwd
#RUN ls -als

RUN go build .

#RUN ls -als

EXPOSE 8080
ENTRYPOINT ["./nats-docker-api"]
