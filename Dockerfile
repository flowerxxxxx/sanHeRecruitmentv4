FROM golang:1.19

MAINTAINER "yanmingyu55@gmail.com"

WORKDIR /home/go/src/sanHeRecruitment

ADD . /home/go/src/sanHeRecruitment

RUN go build main.go

EXPOSE 9090

ENTRYPOINT ["./main"]