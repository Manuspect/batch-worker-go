FROM golang:1.22.5 AS build

WORKDIR /app 

COPY . .

RUN apt-get update && apt-get install make &&  make build


FROM debian:stable-slim

COPY --from=build /app/bin/batch_worker_go /batch_worker_go

CMD ["./batch_worker_go"]