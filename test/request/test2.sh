#!/usr/bin/env bash

set -ex

curl -i -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0MiIsImlhdCI6MTUxNjIzOTAyMn0.-JC0k-LKrqL7fbK4m3k0VjL1X8F06E5Lvb93v1YDjtI" \
    http://localhost:8081/hello
