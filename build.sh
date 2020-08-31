export IMAGE=voteva/ratelimit-operator:v0.0.2
echo $IMAGE

echo "operator-sdk build ..."
operator-sdk build $IMAGE && echo "docker push ..." && docker push $IMAGE