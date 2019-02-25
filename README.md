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

docker-compose exec kafka bash -c "echo 'abc' | kafka-console-producer --request-required-acks 1 --broker-list kafka:9092 --topic heroku_logs"

docker-compose exec kafka bash -c "echo '{\"syslog.pri\":\"190\",\"syslog.timestamp\":\"2019-02-25T16:09:23.296535+00:00\",\"syslog.hostname\":\"host\",\"syslog.appname\":\"app\",\"syslog.procid\":\"sidekiq_pusher.1\",\"message\":\"abc\",\"syslog.facility\":\"local7\",\"syslog.severity\":\"info\",\"logplex.drain_id\":\"d.133fd389-68d7-4d1e-a40d-521da502593a\",\"logplex.frame_id\":\"EBC908AFB0D5433EFE95C18457C0B5F8\"}' | kafka-console-producer --request-required-acks 1 --broker-list kafka:9092 --topic heroku_logs"

```

## Other debug :

```
docker exec kafka-sumo_kafka_1 kafka-topics --create --topic devops --partitions 1 --replication-factor 1 --if-not-exists --zookeeper=zookeeper:2181

docker exec kafka-sumo_kafka_1 kafka-topics --zookeeper=zookeeper:2181 --describe --topic devops

docker exec kafka-sumo_kafka_1 kafka-console-consumer --bootstrap-server kafka:29092 --topic devops --from-beginning
```

## Todo
