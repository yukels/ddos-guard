#!/usr/bin/env bash

set -ex

curl -i -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X2Jsb2NrZWQiLCJpYXQiOjE1MTYyMzkwMjJ9.dWYUv5ErMMIkN62VNXZCG6CTuBAqEm3HoHIeTSz9nD0" \
    http://localhost:8081/hello
