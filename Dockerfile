FROM golang:alpine as builder

ARG GOLANG_NAMESPACE="github.com/mohemohe/mastodon-dynamic-patch"
ENV GOLANG_NAMESPACE="$GOLANG_NAMESPACE"

RUN apk --no-cache add alpine-sdk coreutils git tzdata nodejs upx
RUN cp -f /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
RUN go get -u -v github.com/pwaller/goupx
WORKDIR /go/src/$GOLANG_NAMESPACE
ADD ./go.mod /go/src/$GOLANG_NAMESPACE
ADD ./go.sum /go/src/$GOLANG_NAMESPACE
ENV GO111MODULE=on
RUN go mod download
ADD . /go/src/$GOLANG_NAMESPACE/
RUN go build -ldflags "\
      -w \
      -s \
    " -o /patcher
RUN goupx /patcher

# ====================================================================================

FROM alpine

RUN apk --no-cache add ca-certificates
COPY --from=builder /etc/localtime /etc/localtime
COPY --from=builder /patcher /patcher

EXPOSE 8080
WORKDIR /
CMD ["/patcher"]
