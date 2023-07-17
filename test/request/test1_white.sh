#!/usr/bin/env bash

set -ex

curl -i -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0MV93aGl0ZSIsImlhdCI6MTUxNjIzOTAyMn0.jTgb7vLwG40PzN8orFPevRSL7DOG0M8oqb5oKctmxng" \
    http://localhost:8081/hello
