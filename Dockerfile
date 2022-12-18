FROM golang:1.19.4-alpine as builder


WORKDIR app

ADD src ./

# Install any required modules
RUN go mod download
RUN go build -o server

CMD [ "./server" ]

