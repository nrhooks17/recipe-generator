#!/bin/bash

./start_local_database.sh

sleep 2 # mainly to allow postgres to start

./build_and_start_app.sh

echo "Local setup for recipe generator complete."
