FROM golang:latest 
ADD . /go/src/github.com/dirges/todo
WORKDIR /go/src/github.com/dirges/todo
RUN go build
RUN pwd
RUN ls -ahl
ENTRYPOINT ./todo
