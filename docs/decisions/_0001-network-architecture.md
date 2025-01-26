# Type of Network Architecture

## Context and Problem Statement

The data between clients need to be synchronized. 


## Considered Options

* Server-Client
* Peer-to-Peer

## Decision Outcome

Chosen option: "Peer-to-Peer"

### Consequences

* **Good**:
    * No central server is required
    * Reducing risk of point of faulure

* **Bad**:
    * Each client / peer needs to manage its own state
    * Synchronization and consistency of data can be more complex