# Redis Server

This challenge is to build your own Redis Server.

Redis is an in-memory data structure server, which supports storing strings, hashes, lists, sets, sorted sets and more.

The name Redis reflects the original goal to be a Remote Dictionary Server. **Salvatore Sanfilippo** the creator of Redis 
originally wrote it in just over 300 lines of TCL, you can see that original version in a gist he posted 
[https://gist.github.com/antirez/6ca04dd191bdb82aad9fb241013e88a8](here.)

Since the first version in 2009, Redis has been ported to C and released as open source. It’s also become one of the 
most widely used key-value / NoSQL databases.

# Benchmarking

```
~ > redis-benchmark -t SET,GET -q
WARNING: Could not fetch server CONFIG
SET: 81037.28 requests per second, p50=0.327 msec
GET: 78003.12 requests per second, p50=0.359 msec
```