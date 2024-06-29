#!/usr/bin/bash

docker build -t echo:latest .

docker tag echo:latest hazik1024/echo:latest