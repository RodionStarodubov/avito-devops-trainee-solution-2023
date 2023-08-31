package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Database struct {
	Client  *redis.Client
	Context context.Context
}

var (
	AppAddress    = "localhost:8089"
	RedisAddress  = "localhost:6379"
	RedisUsername = "default"
	RedisPassword = ""
)

func (db *Database) SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	keyValue := KeyValue{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(&keyValue)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if keyValue.Key == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = db.Client.Set(db.Context, keyValue.Key, keyValue.Value, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (db *Database) GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	value, err := db.Client.Get(db.Context, key).Result()
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	keyValue := KeyValue{key, value}

	json, err := json.Marshal(keyValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (db *Database) DelKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	keyValue := KeyValue{}
	err := d.Decode(&keyValue)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if keyValue.Key == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = db.Client.Del(db.Context, keyValue.Key).Err()
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func main() {
	if value, ok := os.LookupEnv("APP_ADDRESS"); ok {
		AppAddress = value
	}

	if value, ok := os.LookupEnv("REDIS_ADDRESS"); ok {
		RedisAddress = value
	}

	if value, ok := os.LookupEnv("REDIS_USERNAME"); ok {
		RedisUsername = value
	}

	if value, ok := os.LookupEnv("REDIS_PASSWORD"); ok {
		RedisPassword = value
	}

	cert, err := tls.LoadX509KeyPair("redis.crt", "redis.key")
	if err != nil {
		log.Fatal("Cannot load redis cert!")
		return
	}

	db := Database{
		redis.NewClient(&redis.Options{
			Addr:     RedisAddress,
			Username: RedisUsername,
			Password: RedisPassword,
			TLSConfig: &tls.Config{
				MinVersion:   tls.VersionTLS12,
				Certificates: []tls.Certificate{cert},
			},
		}),
		context.TODO(),
	}

	http.HandleFunc("/set_key", db.SetKeyHandler)
	http.HandleFunc("/get_key", db.GetKeyHandler)
	http.HandleFunc("/del_key", db.DelKeyHandler)
	http.ListenAndServeTLS(AppAddress, "server.crt", "server.key", nil)
}
