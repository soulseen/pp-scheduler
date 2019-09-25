FROM golang:1.12

ENV CGO_ENABLED=1

RUN apt-get update && apt-get install -y apt-transport-https jq openssl libltdl7 \
    && curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
    && echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list \
    && apt-get update \
    && apt-get install -y kubectl