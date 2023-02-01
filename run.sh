docker load < psyweb_img.tar
docker run --name psyweb_img -e MYSQL_ROOT_PASSWORD=123456 -p 3307:3306 -p 80:8888 -d psyweb_img sh -c 'service mysql start && cd /go/src/jr/psyweb && mysql -uroot -p$MYSQL_ROOT_PASSWORD < db_init.sql && ./psyweb'