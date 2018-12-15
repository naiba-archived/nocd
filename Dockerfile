FROM golang:alpine AS binarybuilder
# Install build deps
RUN apk --no-cache --no-progress add --virtual build-deps build-base git linux-pam-dev
WORKDIR /go/src/github.com/naiba/nocd/
COPY . .
RUN go get -u github.com/tmthrgd/go-bindata/go-bindata \
    && cd cmd/web \
    && go-bindata resource/... \
    && go build -ldflags="-s -w"

FROM alpine:latest
RUN echo http://dl-2.alpinelinux.org/alpine/edge/community/ >> /etc/apk/repositories \
  && apk --no-cache --no-progress add \
    tzdata
# Copy binary to container
WORKDIR /data
COPY --from=binarybuilder /go/src/github.com/naiba/nocd/cmd/web/web ./nocd

# Configure Docker Container
VOLUME ["/data/conf"]
EXPOSE 8000
CMD ["/data/nocd"]
