FROM golang:1.17.8-bullseye

WORKDIR /studier/server

COPY . .

RUN go build .

ENTRYPOINT ["./server"]
