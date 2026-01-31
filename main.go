package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type KVStore struct {
	data   map[string]string
	mu     sync.RWMutex
	conns  []*websocket.Conn
	connMu sync.Mutex
}

func (k *KVStore) Set(key, value string) {
	k.mu.Lock()
	k.data[key] = value
	k.mu.Unlock()
	k.broadcast(key, value)
}

func (k *KVStore) Get(key string) (string, bool) {
	k.mu.RLock()
	v, ok := k.data[key]
	k.mu.RUnlock()
	return v, ok
}

func (k *KVStore) GetAll() map[string]string {
	k.mu.RLock()
	copy := make(map[string]string)
	for k, v := range k.data {
		copy[k] = v
	}
	k.mu.RUnlock()
	return copy
}

func (k *KVStore) broadcast(key, value string) {
	msg := map[string]string{"key": key, "value": value}
	data, _ := json.Marshal(msg)
	k.connMu.Lock()
	for _, conn := range k.conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			// remove conn if error, but for simplicity
		}
	}
	k.connMu.Unlock()
}

func (k *KVStore) addConn(conn *websocket.Conn) {
	k.connMu.Lock()
	k.conns = append(k.conns, conn)
	k.connMu.Unlock()
}

func (k *KVStore) removeConn(conn *websocket.Conn) {
	k.connMu.Lock()
	for i, c := range k.conns {
		if c == conn {
			k.conns = append(k.conns[:i], k.conns[i+1:]...)
			break
		}
	}
	k.connMu.Unlock()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(kv *KVStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		kv.addConn(conn)
		defer kv.removeConn(conn)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}
}

func setHandler(kv *KVStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		if key == "" || value == "" {
			http.Error(w, "missing key or value", 400)
			return
		}
		kv.Set(key, value)
		w.WriteHeader(200)
		fmt.Fprint(w, "ok")
	}
}

func getHandler(kv *KVStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "missing key", 400)
			return
		}
		value, ok := kv.Get(key)
		if !ok {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, value)
	}
}

func getAllHandler(kv *KVStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := kv.GetAll()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all)
	}
}

func main() {
	kv := &KVStore{
		data:  make(map[string]string),
		conns: make([]*websocket.Conn, 0),
	}

	http.HandleFunc("/set", setHandler(kv))
	http.HandleFunc("/get", getHandler(kv))
	http.HandleFunc("/getall", getAllHandler(kv))
	http.HandleFunc("/info-ws", wsHandler(kv))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}