package node

type VirtualNode struct {
	Hash uint64
	Node *Node
}

type Node struct {
	Name string
	Endpoint string
}

