## Build :

```
./docker_build.sh
```

## Run as testing environment with docker-compose :

```
./run.sh
```

## Test

## Debug :

```
docker logs -f kafka-sumo_kafka-sumo_1

docker exec kafka-sumo_kafka_1 kafka-topics --list --zookeeper=zookeeper:2181

docker exec -ti kafka-sumo_kafka_1 bash
kafka-console-producer --broker-list kafka:29092 --topic devops
```

## Other debug :

```
docker exec kafka-sumo_kafka_1 kafka-topics --create --topic devops --partitions 1 --replication-factor 1 --if-not-exists --zookeeper=zookeeper:2181

docker exec kafka-sumo_kafka_1 kafka-topics --zookeeper=zookeeper:2181 --describe --topic devops

docker exec kafka-sumo_kafka_1 kafka-console-consumer --bootstrap-server kafka:29092 --topic devops --from-beginning
```

## Todo
