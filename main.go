package main

import (
	"fmt"

	"github.com/VishalSainani/assignment02bca/blockchain"
)

var bootstrapNodePort string = "6000"

func main() {
	var bootstrapNode = &blockchain.BootstrapNode{}
	blockchain.AddBootstrapNode(bootstrapNode, bootstrapNodePort)
	var node1 = &blockchain.Node{}
	blockchain.AddNode(node1, "6001")
	var node2 = &blockchain.Node{}
	blockchain.AddNode(node2, "6002")
	var node3 = &blockchain.Node{}
	blockchain.AddNode(node3, "6003")
	var node4 = &blockchain.Node{}
	blockchain.AddNode(node4, "6004")
	var node5 = &blockchain.Node{}
	blockchain.AddNode(node5, "6005")
	var node6 = &blockchain.Node{}
	blockchain.AddNode(node6, "6006")
	var node7 = &blockchain.Node{}
	blockchain.AddNode(node7, "6007")
	var node8 = &blockchain.Node{}
	blockchain.AddNode(node8, "6008")
	var node9 = &blockchain.Node{}
	blockchain.AddNode(node9, "6009")
	var node10 = &blockchain.Node{}
	blockchain.AddNode(node10, "6010")
	fmt.Println(bootstrapNodePort, bootstrapNode.PeersList)
	fmt.Println(node1.PortNo, node1.PeersList)
	fmt.Println(node2.PortNo, node2.PeersList)
	fmt.Println(node3.PortNo, node3.PeersList)
	fmt.Println(node4.PortNo, node4.PeersList)
	fmt.Println(node5.PortNo, node5.PeersList)
	fmt.Println(node6.PortNo, node6.PeersList)
	fmt.Println(node7.PortNo, node7.PeersList)
	fmt.Println(node8.PortNo, node8.PeersList)
	fmt.Println(node9.PortNo, node9.PeersList)
	fmt.Println(node10.PortNo, node10.PeersList)
	fmt.Println(bootstrapNode.NodesPortList)

	go blockchain.FloodTransaction(node1, "n1-n2-5")
	go blockchain.FloodTransaction(node2, "n2-n5-1")
	go blockchain.FloodTransaction(node5, "n5-n10-4")
	go blockchain.FloodTransaction(node6, "n6-n3-3")
	for {
	}
}
