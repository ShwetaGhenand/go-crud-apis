FROM golang:1.16-alpine

WORKDIR /user-service

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /userapis

EXPOSE 8081

CMD [ "/userapis" , "server"]
