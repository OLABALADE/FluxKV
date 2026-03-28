package cluster

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sort"
)

type Node struct {
	ID       string
	Address  string
	hashRing *HashRing
}

func NewNode(id string, address string, peers []string) *Node {
	hr := NewHashRing(3)
	peers = append(peers, address)
	sort.Strings(peers)

	for _, peer := range peers {
		hr.AddNode(peer)
	}

	return &Node{
		ID:       id,
		Address:  address,
		hashRing: hr,
	}
}

func (n *Node) ForwardToNode(method, url string, body any) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func (n *Node) GetResponsibleNode(key string) string {
	return n.hashRing.GetNode(key)
}
