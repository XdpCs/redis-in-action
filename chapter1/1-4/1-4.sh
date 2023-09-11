#!/bin/zsh

redis-cli <<EOF
HSET hash-key sub-key1 value1
HSET hash-key sub-key2 value2
HSET hash-key sub-key1 value1
HGETALL hash-key
HDEL hash-key sub-key2
HDEL hash-key sub-key2
HGET hash-key sub-key1
HGETALL hash-key
EOF
