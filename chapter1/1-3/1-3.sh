#!/bin/zsh

redis-cli <<EOF
SADD set-key item
SADD set-key item2
SADD set-key item3
SADD set-key item
SMEMBERS set-key
SISMEMBER set-key item4
SISMEMBER set-key item
SREM set-key item2
SREM set-key item2
SMEMBERS set-key
EOF
