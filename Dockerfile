FROM golang

WORKDIR /go/src/beancount-uploader

ADD . /go/src/beancount-uploader/

RUN go get github.com/tools/godep

RUN godep restore

RUN go install beancount-uploader

ENV SENDER abc@gmail.com
ENV GIN_MODE release
ENV AWS_ACCESS_KEY_ID key_id
ENV AWS_SECRET_ACCESS_KEY access_key
ENV BEANCOUNT_BUCKET bucket
ENV AWS_REGION eu-west-1

EXPOSE 8080

ENTRYPOINT ["/go/bin/beancount-uploader"]