FROM golang:1.18

RUN mkdir /api

ADD . /api
WORKDIR /api

RUN go mod tidy
RUN go build -o api .
RUN chmod +x /api/api

ENTRYPOINT ["/api/api"]