FROM golang:1.20-bullseye

ENV ROOT=/go/src

WORKDIR ${ROOT}

RUN apt update \
    && apt clean \
    && rm -r /var/lib/apt/lists/*

RUN apt install git && \
    apt install curl

COPY go.mod go.sum ./

RUN /bin/sh -c /bin/sh -c go mod download

RUN /bin/sh -c /bin/sh go get -u github.com/go-sql-driver/mysql

RUN go install github.com/cosmtrek/air@v1.29.0

COPY . .

EXPOSE 8080

# CMD ["air", "-c", ".air.toml"]

# CMD ["go", "run", "cmd/main.go"]
