FROM golang:alpine as builder

COPY . /go/src/github.com/Evaneos/concourse-gcloudsql-resource
RUN go build -o /assets/check github.com/Evaneos/concourse-gcloudsql-resource/check
RUN go build -o /assets/check github.com/Evaneos/concourse-gcloudsql-resource/in

FROM google/cloud-sdk:alpine

COPY --from=builder /assets/check /assets/check
COPY --from=builder /assets/check /assets/in
