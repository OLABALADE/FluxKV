package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/OLABALADE/FluxKV/internal/cluster"
	"github.com/OLABALADE/FluxKV/internal/store"
)

type Handler struct {
	store   store.Store
	cluster *cluster.Node
}

func NewHandler(s store.Store, c *cluster.Node) *Handler {
	return &Handler{store: s, cluster: c}
}

type Request struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	req := Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	targetNode := h.cluster.GetResponsibleNode(req.Key)
	// If not target Node
	if targetNode != h.cluster.Address {
		h.FowardRequest(w, "put", targetNode, "", req)
		return
	}

	log.Println("I'm responsible I", targetNode)
	if err := h.store.Put(req.Key, req.Value); err != nil {
		log.Println("ERROR: ", err)
		http.Error(w, "Failed to add key to store", http.StatusInternalServerError)
		return
	}

	h.cluster.Replicate("POST", req.Key, req.Value)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	// targetNode := h.cluster.GetResponsibleNode(key)
	// // // If not target Node
	// // if targetNode != h.cluster.Address {
	// // 	h.FowardRequest(w, "get", targetNode, key, nil)
	// // 	return
	// // }
	// //
	log.Println("I'm responsible I", h.cluster.Address)
	val, err := h.store.Get(key)
	if err != nil {
		log.Println("ERROR: ", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"value": val,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	// targetNode := h.cluster.GetResponsibleNode(key)
	//If not target Node
	// if targetNode != h.cluster.Address {
	// 	h.FowardRequest(w, "delete", targetNode, key, nil)
	// 	return
	// }

	log.Println("I'm responsible I", h.cluster.Address)
	if err := h.store.Delete(key); err != nil {
		log.Println("ERROR: ", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	h.cluster.Replicate("DELETE", key, "")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) FowardRequest(w http.ResponseWriter, method, targetNode, key string, body any) {
	log.Println("I'm not responsible ask", targetNode)
	method = strings.ToLower(method)

	url := ""

	urlStr := "http://%s/"
	if method == "get" || method == "delete" {
		urlStr = urlStr + method + "?key=%s"
		url = fmt.Sprintf(urlStr, targetNode, key)
	} else {
		url = fmt.Sprintf(urlStr, targetNode) + "put"
	}

	resp, err := h.cluster.ForwardToNode(method, url, body)

	if err != nil {
		log.Println("ERROR: Failed to forward request", err)
		http.Error(w, "Failed to forward request", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	if method == "get" {
		io.Copy(w, resp.Body)
	}
}

func (h *Handler) ReplicatePut(w http.ResponseWriter, r *http.Request) {
	req := Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if err := h.store.Put(req.Key, req.Value); err != nil {
		log.Println("ERROR: ", err)
		http.Error(w, "Failed to add key to store", http.StatusInternalServerError)
		return
	}

	log.Println("INFO: Replicated data into Store")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ReplicateDelete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if err := h.store.Delete(key); err != nil {
		log.Println("ERROR: ", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	log.Println("INFO: Removed replica data from Store")
	w.WriteHeader(http.StatusOK)
}
