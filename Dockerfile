FROM golang:1.19-buster

ENV GOPATH=/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./ ./

RUN go build -o backend ./server/cmd/v1/main.go

CMD ["./backend"]