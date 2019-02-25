#!/bin/bash

docker login && docker build -t vietwow/kafka-sumo . && docker push vietwow/kafka-sumo
