FROM golang:1.13

ENV GO111MODULE=on
ENV PORT=8080
WORKDIR /app/server
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o main
ENTRYPOINT ["./main"]
CMD ["-webmode=true"]