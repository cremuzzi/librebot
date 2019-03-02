FROM golang:1.12.0-alpine3.9

WORKDIR /go/src/app
COPY . .

RUN apk add --no-cache git 
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]
