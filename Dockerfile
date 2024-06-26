FROM ubuntu:latest
LABEL authors="oloja"

ENTRYPOINT ["top", "-b"]