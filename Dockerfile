FROM golang:1.23-alpine AS builder

RUN apk --no-cache add make gcc musl-dev

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o ./bin/authservice cmd/authservice/main.go


FROM alpine

COPY --from=builder /app/bin/authservice /
COPY configs/config.yaml configs/config.yaml

CMD ["/authservice"]