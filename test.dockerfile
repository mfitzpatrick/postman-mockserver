FROM golang:alpine3.12

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /test
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
CMD ["sh", "+e", "-c", "for m in cmd common postman ; do cd /test/$m ; go test ; done"]

