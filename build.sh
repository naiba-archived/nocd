#!/usr/bin/env bash
cd cmd/web &&\
    go-bindata resource/... &&\
    go build &&\
    mv web ~/Downloads/cd/