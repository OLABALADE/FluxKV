package cluster

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var replicateFactor = 2

type ReplicationRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (n *Node) Replicate(method, key, value string) {
	//Get Peers from HashRing
	seen := map[string]struct{}{}
	for _, vnode := range n.hashRing.nodes {
		if len(seen) >= replicateFactor-1 {
			break
		}

		node, _ := n.hashRing.vNodeMap[vnode]

		if node == n.Address {
			continue
		}

		_, ok := seen[node]

		if ok {
			continue
		}

		seen[node] = struct{}{}

		//Send to Peers
		go func(p string) {
			// var b []byte = nil
			var body any = nil

			if method == "POST" {
				body = ReplicationRequest{
					Key:   key,
					Value: value,
				}

			}

			url := ""
			if method == "DELETE" {
				url = fmt.Sprintf("http://%s/replicate?key=%s", p, key)
			} else {
				url = fmt.Sprintf("http://%s/replicate", p)
			}

			resp, err := n.ForwardToNode(method, url, body)

			if err != nil {
				log.Println("Replication network error to", p, err)
				return
			}
			defer resp.Body.Close()

			bodyBytes, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != http.StatusOK {
				log.Printf("Replication failed to %s | status=%d | body=%s\n",
					p,
					resp.StatusCode,
					string(bodyBytes),
				)
				return
			}

			log.Printf("Replicated to %s | status=%d\n", p, resp.StatusCode)

		}(node)
	}

}
