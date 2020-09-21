export APP_NAME=ratelimit-operator
export VERSION=v0.0.0
export IMAGE=voteva/$APP_NAME:$VERSION

echo $IMAGE

#go test ./pkg/... &&
mkdir -p build/_output/bin &&
go build -o ./build/_output/bin/$APP_NAME -v ./cmd/manager &&
echo "docker build ..." &&
docker build -t $IMAGE .
echo "docker push ..." &&
docker push $IMAGE
