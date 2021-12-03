#!/usr/bin/env sh

export MQ_HOST=localhost
export MQ_PORT=5672
export MQ_USER=user
export MQ_PASSWORD=password

printf "Starting RabbitMQ container"
docker run -d --rm --name mlmodelscope_rabbitmq-integration -p 5672:5672 --env RABBITMQ_DEFAULT_USER="${MQ_USER}" --env RABBITMQ_DEFAULT_PASS="${MQ_PASSWORD}" rabbitmq:3-alpine > /dev/null

while [ -z "`docker logs mlmodelscope_rabbitmq-integration | grep "started TCP listener"`" ]
do
    printf "."
    sleep 1
done

echo "\nRunning integration tests..."
go clean -testcache && go test --tags integration ./...

echo "Cleaning up RabbitMQ container..."
docker stop mlmodelscope_rabbitmq-integration > /dev/null