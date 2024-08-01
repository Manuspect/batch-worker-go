FROM golang:1.22.5 AS build

WORKDIR /app 

COPY . .

RUN apt-get update && apt-get install make &&  make build


FROM jrottenberg/ffmpeg:4-ubuntu

COPY --from=build /app/bin/batch_worker_go /batch_worker_go

CMD ["./batch_worker_go"]