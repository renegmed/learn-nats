FROM golang:1.14.0-alpine

WORKDIR src/worker

COPY . .
 
RUN go mod vendor
RUN go build -o worker .

#RUN ls -als

EXPOSE 8181
ENTRYPOINT ["./worker"]
