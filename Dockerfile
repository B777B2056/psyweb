FROM golang:1.19-buster AS golang_builder
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /go/src/jr/psyweb
COPY . /go/src/jr/psyweb
RUN go build -o psyweb main.go

FROM python:3.7-buster
WORKDIR /go/src/jr/psyweb
COPY --from=golang_builder /go/src/jr/psyweb/psyweb .
COPY ./db_init.sql ./db_init.sql
COPY ./mysql-apt-config_0.8.24-1_all.deb ./mysql-apt-config_0.8.24-1_all.deb
COPY ./configuration ./configuration
COPY ./web/views ./web/views
COPY ./deeplearning ./deeplearning
RUN rm -f ./web/views/views.go \
&& apt-get update && apt install -y mariadb-server \
&& pip3 install -r deeplearning/requirements.txt -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com --user \
&& pip3 install torch --extra-index-url https://download.pytorch.org/whl/cpu -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com --user

EXPOSE 8888
# ENTRYPOINT ["./psyweb"]
