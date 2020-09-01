export IMAGE=d4rkest/ratelimit-operator:v0.0.2
echo $IMAGE

echo "operator-sdk build ..."
go test ./pkg/... && operator-sdk build $IMAGE && echo "docker push ..." && docker push $IMAGE