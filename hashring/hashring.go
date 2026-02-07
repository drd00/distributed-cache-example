package hashring

import (
	"fmt"
	"hash/fnv"
	"sort"

	"github.com/drd00/distributed-cache-example/node"
)

const (
	HashesPerKey int = 150
)

type HashRing interface {
	Init()
	GetNode(key string) (*node.Node, error)
	getIndexForKey(key string) (int, error)
	assignVirtualNodes() error
	hash(key string) uint64
}

type HashRingImpl struct {
	virtualNodes []*node.VirtualNode
	nodes []*node.Node
}

func NewHashRing(nodes []*node.Node) HashRing {
	hr := &HashRingImpl{
		nodes: nodes,
	}
	hr.Init()

	return hr
}

func (hr *HashRingImpl) Init() {
	nNodes := len(hr.nodes)
	hrSize := nNodes * HashesPerKey
	hr.virtualNodes = make([]*node.VirtualNode, 0, hrSize)
}

func (hr *HashRingImpl) getIndexForKey(key string) (int, error) {
	if len(hr.virtualNodes) == 0 {
		return -1, fmt.Errorf(
			"cannot get index for key %s: `virtualNodes` member variable has length 0",
			key,
		)
	}

	// Hash the key value to get a uint64 value
	keyHash := hr.hash(key)

	// Binary search to find first value which is >= `keyHash`
	idx := sort.Search(len(hr.virtualNodes), func(i int) bool {
		return hr.virtualNodes[i].Hash >= keyHash
	})

	if idx >= len(hr.virtualNodes) {
		idx = 0
	}

	return idx, nil
}

func (hr *HashRingImpl) GetNode(key string) (*node.Node, error) {
	idx, err := hr.getIndexForKey(key)
	if err != nil {
		return nil, err
	}

	return hr.virtualNodes[idx].Node, nil
}

func (hr *HashRingImpl) assignVirtualNodes() error {
	if hr.virtualNodes == nil {
		return fmt.Errorf("cannot add values to virtualNodes: member virtualNodes is nil")
	}

	// Hash each node `HashesPerKey` times with a different salt
	// And append to `hr.virtualNodes` slice
	for _, hrNode := range hr.nodes {
		for i := range HashesPerKey {
			virtualNode := &node.VirtualNode{
				Hash: hr.hash(fmt.Sprintf("%s-%d", hrNode.Name, i)),
				Node: hrNode,
			}

			hr.virtualNodes = append(hr.virtualNodes, virtualNode)
		}
	}

	// Sort the `virtualNodes` slice in ascending order by hash value
	sort.Slice(
		hr.virtualNodes,
		func (i, j int) bool {
			return hr.virtualNodes[i].Hash < hr.virtualNodes[j].Hash
		},
	)

	return nil
}

func (hr *HashRingImpl) hash(value string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(value))
	return h.Sum64()
}

