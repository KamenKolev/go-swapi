FROM golang:1.18

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./src/*.go ./

RUN go build -o /swapi-go

EXPOSE 8080

CMD [ "/swapi-go" ]
