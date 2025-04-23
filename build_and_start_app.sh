#!/bin/bash

buildAppCommand="sudo go build -o recipe_generator ./cmd/server"

echo "Building recipe generator backend code with the following command: ${buildAppCommand}..."

go build -o recipe_generator ./cmd/server

echo "App built successfully."

echo "Starting recipe generator backend server..."

./recipe_generator
