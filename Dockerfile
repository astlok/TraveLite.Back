FROM golang:1.16-alpine

WORKDIR /app

COPY . /app

RUN go build -o /bin/trailite -v ./cmd/

COPY ./config/config.yaml /config/

EXPOSE 8080

CMD ["trailite log stdout config /config/config.yaml"]