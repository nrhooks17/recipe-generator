#!/bin/bash

startDockerCommand="sudo systemctl start docker"
echo "Starting docker process by executing '${startDockerCommand}'"
${startDockerCommand}

startPostgresCommand="docker container start postgres-db"
echo "Starting postgres container with command '${startPostgresCommand}'"
${startPostgresCommand}
