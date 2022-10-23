FROM golang:alpine
RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN go build -o benchmark ./cmd/benchmark/benchmark.go
CMD ["./benchmark"]