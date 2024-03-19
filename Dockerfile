FROM golang:latest

LABEL maintainer="Nikita Shlyakhtin <nikitashliahtin@mail.ru>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build/fc

EXPOSE 8080

CMD ["./bin/floodcontrol"]