#!/usr/bin/env bash

set -ex
curl -i -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0MSIsImlhdCI6MTUxNjIzOTAyMn0.gLvATAsFAI9qurYyDkuXPdHM2JZg6al0p1N21dKbgtE" \
    http://localhost:8081/hello