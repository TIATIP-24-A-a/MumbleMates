# Peer Discovery

## Context and Problem Statement

Peers currently need to know their effective addresds to create a connection.
Peers need a way to easily discover each other to establish a communication. 

## Considered Options

* Rendezvous
* mDNS (Multicast DNS)
* DHT (Distributed Hash Table)

## Decision Outcome

Chosen option: "mDNS (Multicast DNS)", because it works well for local network environments and is already support by the libp2p library.

### Consequences

* **Good**:
    * Only for local network
    * Supported by libp2p with zero configuration
    * Works for low number of peers

* **Bad**:
    * Limited to local network only
    * May not scale well on a larger network
    * Verbose when there are a lot of peers