FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./src/*.go ./

RUN go build -o /app

EXPOSE 8080

CMD [ "/app" ]