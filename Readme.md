# Distributed KV Store

# Concept
Multiple kv server, communicating with each other, when a server store gets updated, other store for servers also gets updated. (Replication)
If a server stops, it will be replaced by other server.

# Todo
1. Protocol Design
  - Binary Protocol

2. Parser

3. Node Design

4. Pool Design (as of now it will consist of 1 write node and 1 read node)

