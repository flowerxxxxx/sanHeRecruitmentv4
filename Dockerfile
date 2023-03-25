FROM golang:1.19

MAINTAINER "yanmingyu55@gmail.com"

# 主文件
WORKDIR /home/go/src/sanHeRecruitment

ADD . /home/go/src/sanHeRecruitment

# 静态文件夹
# WORKDIR /home/sanheRecPic

COPY . /home/sanheRecPic

# 备份文件夹
# WORKDIR /home/sanheRecBackup

COPY . /home/sanheRecBackup

RUN go build main.go

EXPOSE 9090

ENTRYPOINT ["./main"]