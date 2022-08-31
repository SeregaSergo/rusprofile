FROM golang:1.19-buster
LABEL maintainer="khattu.s@mail.ru"

COPY ./ /app
WORKDIR /app

RUN apt-get update
RUN apt-get install -y protobuf-compiler

# build go app
RUN make

CMD ["./server"]