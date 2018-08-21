FROM golang:alpine as builder

RUN apk add --no-cache git

COPY . /go/src/github.com/Evaneos/concourse-gcloudsql-resource
RUN go get github.com/google/uuid
RUN go build -o /assets/check github.com/Evaneos/concourse-gcloudsql-resource/check
RUN go build -o /assets/in github.com/Evaneos/concourse-gcloudsql-resource/in

FROM google/cloud-sdk:alpine

COPY --from=builder /assets/check /opt/resource/check
COPY --from=builder /assets/in /opt/resource/in
