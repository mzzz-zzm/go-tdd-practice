#make sure to specify the same Go version as the one in the go.mod file
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .
ARG BIN_TO_BUILD
RUN go build -o svr cmd/${BIN_TO_BUILD}/*.go
ARG PORT_TO_EXPOSE
EXPOSE ${PORT_TO_EXPOSE}
CMD [ "./svr" ]