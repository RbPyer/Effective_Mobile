FROM golang:1.23.1 AS builder

WORKDIR /my_app

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile ./
COPY . .

RUN apt-get update && apt-get install -y make

RUN make scratch_build


FROM alpine:latest
WORKDIR alpine_app

COPY --from=builder /my_app/bin/app /bin/app
COPY --from=builder /my_app/config config


ENTRYPOINT ["/bin/app"]