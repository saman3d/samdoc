FROM alpine:3.17.0

ARG PUID=1000
ARG PGID=1000

ENV TZ=Asia/Tehran
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apk update
RUN apk add --no-cache libreoffice-writer 
RUN apk add ttf-freefont 
COPY ./fonts/ /usr/share/fonts

ENTRYPOINT ["/usr/bin/libreoffice"]
