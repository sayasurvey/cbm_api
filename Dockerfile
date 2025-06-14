FROM golang:1.23-bullseye

WORKDIR /app

RUN apt update && \
    apt install -y git curl && \
    apt clean && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

RUN go get -u gorm.io/gorm
RUN go get -u gorm.io/driver/postgres

RUN go install github.com/cosmtrek/air@v1.29.0
ENV PATH="/go/bin:${PATH}"

COPY . .

ENV PORT=8080
EXPOSE 8080

# 開発用途
CMD ["air", "-c", ".air.toml"]

# 本番用途
CMD ["go", "run", "cmd/main.go"]
