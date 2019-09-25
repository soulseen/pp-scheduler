#!/bin/bash
set -e

function cleanup(){
    result=$?
    echo "Cleaning"
    kubectl delete ns $TEST_NS
    exit $result
}

dest="./deploy/ks-scheduler.yaml"
tag=test-e2e
IMG=zhuxiaoyang/ks-scheduler:$tag
TEST_NS=scheduler-test

trap cleanup EXIT SIGINT SIGQUIT
docker build -f Dockerfile -t ${IMG} .
#docker push $IMG

kubectl create ns  $TEST_NS
kubectl create -f $dest

export TEST_NS

go test -mod=vendor -v ./test/e2e/