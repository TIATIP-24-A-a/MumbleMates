package internal

import mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"

// Starts the mDNS discovery service
func setupMDNSDiscovery(chatNode *ChatNode) error {
	service := mdns.NewMdnsService(chatNode.Node, SERVICE_TAG, chatNode)
	return service.Start()
}
