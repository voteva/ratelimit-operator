export IMAGE=voteva/ratelimit-operator:v0.0.1
echo $IMAGE

echo "operator-sdk build ..."
go test ./pkg/... && operator-sdk build $IMAGE && echo "docker push ..." && docker push $IMAGE