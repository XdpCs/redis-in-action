#!/bin/zsh

redis-cli <<EOF
SET hello world
GET hello
DEL hello
GET hello
EOF
