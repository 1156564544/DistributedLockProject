# Distributed Lock Design
- **Design a simple consensus system, which satisfy the following requirements
- **Support multiple clients to preempt/release a distributed lock, and check the owner of a distributed lock
- **To ensure the data consistency of the system, the follower servers send all preempt/release requests to the leader server
- **To check the owner of a distributed lock, the follower server accesses its map directly and sends the results to the clients
- **When the leader server handling preempt/release requests
- **In this system, all clients provide preempt/release/check distributed lock interface
