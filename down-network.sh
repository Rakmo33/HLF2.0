docker-compose -f ./artifacts/channel/create-certificate-with-ca/docker-compose.yaml stop

echo "waiting..."
sleep 2

docker-compose -f ./artifacts/docker-compose.yaml stop