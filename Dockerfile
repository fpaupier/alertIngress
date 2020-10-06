FROM golang:1.15.2-buster AS builder

WORKDIR /go/src/app

COPY . .

RUN go get -v github.com/confluentinc/confluent-kafka-go/kafka
RUN go get -v github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres
RUN go get -v github.com/golang/protobuf/proto

# Build the application
RUN go build -o main .

# Command to run when starting the container
CMD ["/go/src/app/main"]