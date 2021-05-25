FROM golang:1.15

RUN mkdir /app
WORKDIR /app

#COPY go.mod .
#COPY go.sum .


COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /myapp

ENTRYPOINT ["/myapp"]