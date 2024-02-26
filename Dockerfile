FROM golang:latest AS baseImage

ENV GOPROXY=https://goproxy.cn

RUN rm -rf /var/www
RUN mkdir -p /var/www/
COPY .. /var/www/
WORKDIR /var/www/project-api
RUN mkdir -p bin
RUN ls -al
RUN go mod download && go build -o bin/main
WORKDIR /var/www/project-api/bin
RUN cp ../config . -rf

FROM alpine:latest
ENV TZ=Asia/Shanghai
COPY --from=baseImage /var/www/project-api/bin /app
WORKDIR /app
CMD ["./main"]