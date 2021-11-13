FROM golang:1.16-alpine

WORKDIR /src 

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /grpc-apis

EXPOSE 8081

CMD [ "/grpc-apis" , "server"]
