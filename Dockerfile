FROM golang:1.14

RUN mkdir /dummy
WORKDIR /dummy

COPY go.mod .

RUN go mod download

ADD src src

COPY run.sh .
RUN chmod +x run.sh

EXPOSE 8888
ENTRYPOINT ["./run.sh"]