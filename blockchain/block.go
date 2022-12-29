package blockchain

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"time"
)

var bootstrapNodePort string = "6000"

type Node struct {
	PortNo           string
	PeersList        []string
	TransactionsList []string
}

type BootstrapNode struct {
	Node
	NodesPortList []string
}

func AddBootstrapNode(node *BootstrapNode, PortNo string) {
	node.NodesPortList = append(node.NodesPortList, PortNo)
	node.Node.PortNo = PortNo
	go AsBootstrapServer(node)
}

func AddNode(node *Node, PortNo string) {
	node.PortNo = PortNo
	go AsServer(node)
	JoinNetwork(node)
}

func HandleJoinRequest(c net.Conn, node *BootstrapNode, senderPort string) {
	rand.Seed(time.Now().UnixNano())
	var randOrderNodes = rand.Perm(len(node.NodesPortList) - 0)
	var randNodesToSend []string
	for i := 0; i < int(math.Log2(float64(len(randOrderNodes))))+1; i++ {
		randNodesToSend = append(randNodesToSend, node.NodesPortList[randOrderNodes[i]])
	}

	enc := gob.NewEncoder(c)
	enc.Encode(randNodesToSend)
	node.NodesPortList = append(node.NodesPortList, senderPort)
}

func JoinNetwork(node *Node) {
	c, err := net.Dial("tcp", ":"+bootstrapNodePort)
	if err != nil {
		log.Fatal(err)
	}
	var n [16]byte
	copy(n[:], "Join Request"+node.PortNo)
	c.Write(n[0:16])

	dec := gob.NewDecoder(c)
	err = dec.Decode(&node.PeersList)

	ConnectWithNodes(node)
}

func AsServer(node *Node) {
	ln, err := net.Listen("tcp", ":"+node.PortNo)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go HandleConnection(conn, node)
	}
}

func AsBootstrapServer(node *BootstrapNode) {
	ln, err := net.Listen("tcp", ":"+node.PortNo)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go HandleBootstrapConnection(conn, node)
	}
}

func HandleConnection(c net.Conn, node *Node) {
	buf := make([]byte, 16)
	n, err := c.Read(buf)
	if err != nil || n == 0 {
		c.Close()
	}
	if string(buf[0:1]) == "n" {
		var rcvdTransaction string = string(buf[0:n])
		if !Contains(node.TransactionsList, rcvdTransaction) {
			FloodTransaction(node, rcvdTransaction)
		}
	} else if n == 4 {
		node.PeersList = append(node.PeersList, string(buf[0:4]))
	}
}

func HandleBootstrapConnection(c net.Conn, node *BootstrapNode) {
	buf := make([]byte, 16)
	n, err := c.Read(buf)
	if err != nil || n == 0 {
		c.Close()
	}
	if string(buf[0:12]) == "Join Request" {
		HandleJoinRequest(c, node, string(buf[12:16]))
	} else if string(buf[0:1]) == "n" {
		var rcvdTransaction string = string(buf[0:n])
		if !Contains(node.TransactionsList, rcvdTransaction) {
			FloodBootstrapTransaction(node, rcvdTransaction)
		}
	} else if n == 4 {
		node.PeersList = append(node.PeersList, string(buf[0:4]))
	}
}

func ConnectWithNodes(node *Node) {
	for _, i := range node.PeersList {
		c, err := net.Dial("tcp", ":"+i)
		if err != nil {
			log.Fatal(err)
		}
		var n [4]byte
		copy(n[:], node.PortNo)
		c.Write(n[0:4])
	}
}

func FloodTransaction(node *Node, transaction string) {
	node.TransactionsList = append(node.TransactionsList, transaction)
	fmt.Println(node.PortNo, transaction)
	for _, i := range node.PeersList {
		c, err := net.Dial("tcp", ":"+i)
		if err != nil {
			log.Fatal(err)
		}
		n := make([]byte, len(transaction))
		copy(n[:], transaction)
		c.Write(n[0:len(transaction)])
	}
}

func FloodBootstrapTransaction(node *BootstrapNode, transaction string) {
	node.TransactionsList = append(node.TransactionsList, transaction)
	fmt.Println(node.PortNo, transaction)
	for _, i := range node.PeersList {
		c, err := net.Dial("tcp", ":"+i)
		if err != nil {
			log.Fatal(err)
		}
		n := make([]byte, len(transaction))
		copy(n[:], transaction)
		c.Write(n[0:len(transaction)])
	}
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
