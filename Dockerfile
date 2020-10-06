FROM golang:1.15.2-buster AS builder

WORKDIR /go/src/app

COPY . .

RUN go get -d -v

# Build the application
RUN go build -o main .

# Command to run when starting the container
CMD ["/go/src/app/main"]