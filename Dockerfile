FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go build -o main .

FROM alpine

RUN mkdir /lift
COPY --from=builder /build/main /lift/

COPY ./scripts/start.sh /lift/
RUN chmod +x /lift/start.sh

COPY ./.env /lift/
COPY ./templates/ /lift/templates/
COPY ./static/keys/lift.pub /lift/.ssh/id_rsa.pub
COPY ./static/keys/lift /lift/.ssh/id_rsa

WORKDIR /lift

ARG db_user
ARG db_password
ARG db_host

ENV DB_USER=${db_user}
ENV DB_PASSWORD=${db_password}
ENV DB_HOST=${db_host}

CMD ./start.sh --db_user ${DB_USER} --db_password ${DB_PASSWORD} --db_host ${DB_HOST}