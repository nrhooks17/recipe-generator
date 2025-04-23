#!/bin/bash


echo "Stopping recipe generator backend server..."

ps -aux | grep recipe_generator | grep -v grep | awk '{print $2}' | xargs kill -9

echo "Recipe generator backend server stopped."

echo "Stopping postgres container..."

docker container stop postgres-db

echo "Postgres container stopped."

