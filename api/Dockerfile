FROM golang:1.14.0-alpine

WORKDIR src/api

COPY . .
 
RUN go mod vendor
RUN go build -o api .

#RUN ls -als

EXPOSE 8080
ENTRYPOINT ["./api"]
