FROM mysql/mysql-server:latest
COPY ./migrations/init.sql /docker-entrypoint-initdb.d/