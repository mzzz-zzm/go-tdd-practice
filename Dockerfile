#make sure to specify the same Go version as the one in the go.mod file
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o svr cmd/httpserver/*.go
EXPOSE 8080
CMD [ "./svr" ]