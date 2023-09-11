#!/bin/zsh

redis-cli <<EOF
RPUSH list-key item
RPUSH list-key item2
RPUSH list-key item
LRANGE list-key 0 -1
LINDEX list-key 1
LPOP list-key
LRANGE list-key 0 -1
EOF
