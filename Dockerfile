FROM golang:1.23-bullseye

ENV ROOT=/go/src
WORKDIR ${ROOT}

RUN apt update && \
    apt install -y git curl ca-certificates gnupg && \
    curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
    apt install -y nodejs && \
    apt clean && \
    rm -rf /var/lib/apt/lists/*

RUN npm install -g vercel

COPY go.mod go.sum ./
RUN go mod download


RUN go install github.com/cosmtrek/air@v1.29.0

COPY . .

EXPOSE 8080
