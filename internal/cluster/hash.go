package cluster

import (
	"fmt"
	"hash/fnv"
	"slices"
	"sort"
)

type HashRing struct {
	nodes    []uint32
	vNodeMap map[uint32]string
	vNodes   int
}

func NewHashRing(vnodes int) *HashRing {
	return &HashRing{
		vNodeMap: make(map[uint32]string),
		vNodes:   vnodes,
	}
}

func (h *HashRing) AddNode(node string) {
	for i := 0; i < h.vNodes; i++ {
		vnodeKey := fmt.Sprintf("%s_%d", node, i)
		hash := h.generateHash(vnodeKey)
		h.nodes = append(h.nodes, hash)
		h.vNodeMap[hash] = node
	}
	slices.Sort(h.nodes)
}

func (h *HashRing) GetNode(key string) string {

	hash := h.generateHash(key)

	idx := sort.Search(len(h.nodes), func(i int) bool {
		return h.nodes[i] >= hash
	})

	if idx == len(h.nodes) {
		idx = 0
	}

	return h.vNodeMap[h.nodes[idx]]
}

func (h *HashRing) generateHash(key string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return hash.Sum32()
}
