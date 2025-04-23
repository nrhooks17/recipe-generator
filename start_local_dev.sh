#!/bin/bash

startDockerCommand="sudo systemctl start docker."
echo "Starting docker process by executing '${startDockerCommand}'"

sudo systemctl start docker

startPostgresCommand="docker container start postgres-db."
echo "Starting postgres container with command '${startPostgresCommand}'"

docker container start postgres-db

sleep 2 # mainly to allow postgres to start

./build_and_start_app.sh

echo "Local setup for recipe generator complete."
