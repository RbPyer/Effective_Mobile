FROM golang:1.23.1 AS builder

WORKDIR /my_app

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile ./
COPY . .

RUN apt-get update && apt-get install -y make

RUN make migrator_build
RUN make scratch_build


FROM alpine:latest

COPY --from=builder /my_app/migrations /app/migrations
COPY --from=builder /my_app/bin/app /app/bin/app
COPY --from=builder /my_app/config /app/config
COPY --from=builder /my_app/bin/migrator /app/bin/migrator


RUN chmod +x /app/bin/app
RUN chmod +x /app/bin/migrator