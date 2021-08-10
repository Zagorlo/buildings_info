FROM golang:latest AS tmp
ENV CGO_ENABLED 0

COPY . /go/src/buildings_info
WORKDIR /go/src/buildings_info

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

CMD ["/go/src/buildings_info/main"]
EXPOSE 801