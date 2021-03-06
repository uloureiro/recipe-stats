FROM golang:1.15

WORKDIR /go/src/recipe-stats
COPY . .
ADD sample_data.tar.gz .

RUN go get -d -v ./...
RUN go install -i -v ./...
RUN go build

CMD ["/bin/sh"]