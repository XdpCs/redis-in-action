#!/bin/zsh

redis-cli <<EOF
ZADD zset-key 728 member1
ZADD zset-key 982 member0
ZADD zset-key 982 member0
ZRANGE zset-key 0 -1 WITHSCORES
ZRANGEBYSCORE zset-key 0 800 WITHSCORES
ZREM zset-key member1
ZREM zset-key member1
ZRANGE zset-key 0 -1 WITHSCORES
EOF
