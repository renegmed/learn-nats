FROM golang:1.14.0-alpine

WORKDIR src/nats-docker-worker-2

COPY . .

#RUN pwd
#RUN ls -als

RUN go build .

#RUN ls -als

EXPOSE 8181
ENTRYPOINT ["./nats-docker-worker-2"]
