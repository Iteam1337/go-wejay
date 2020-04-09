FROM golang:1.13-buster

WORKDIR /go/src/wejay

COPY . .

RUN make release

CMD "release/wejay/bin"
