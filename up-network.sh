docker-compose -f ./artifacts/channel/create-certificate-with-ca/docker-compose.yaml up -d

echo "waiting..."
sleep 2

docker-compose -f ./artifacts/docker-compose.yaml up -d

sleep 5

./createChannel.sh

sleep 5

./deployChaincode.sh
