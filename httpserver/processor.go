package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func deleteValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		log.Println("Only DELETE supported for this url!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client := redis.NewClient(&redis.Options{
                Addr:     "localhost:6379",
                Password: "", // no password set
                DB:               0,  // use default DB
        })

        ctx := context.Background()
	keyMap := r.URL.Query()
	for _, keys := range keyMap {
		for _, key := range keys {
			//get all values associated with key
			_, err1 := client.Del(ctx, key).Result()
			if err1 != nil {
				log.Println(err1.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Println("Key " + key + " deleted!")
		}	
	}
	w.WriteHeader(http.StatusOK)
	log.Println("Deleting values done!")
}

func insertValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut{
		log.Println("Only POST and PUT supported for this url!")
		w.WriteHeader(http.StatusBadRequest)
                return
        }
	client := redis.NewClient(&redis.Options{
                Addr:     "localhost:6379",
                Password: "", // no password set
                DB:               0,  // use default DB
        })

        ctx := context.Background()
	segmentMap := r.URL.Query()
	ttlPresent := false
	ttl := 0
	//make sure no segments are longer than 32 bytes and ttl exists
	for k, _ := range segmentMap {
		if len(k) > 32 {
			log.Println("Segment(s) length(s) too long! Must be 32 bytes or under.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if k == "ttl" {
			//if there is more than one ttl value in url
			if len(segmentMap[k]) > 1 {
				log.Println("URL invalid! Can only have one ttl.")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			ttlPresent = true
			err := errors.New("")
			ttl, err = strconv.Atoi(segmentMap[k][0])
			if err != nil || ttl < -1 {
				log.Println("TTL field is invalid! Must be an integer -1 or greater.")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		if(r.Method == http.MethodPut) {
			client.Del(ctx, k)
		}
	}
        for k, v := range segmentMap {
		if k == "ttl" {
			continue
		}
		//push all values belonging to a key to a list
		for _, val := range v {
			_, errMsg := client.RPush(ctx, k, val).Result()
			if errMsg != nil {
				log.Println(errMsg.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		//set ttl
		if ttlPresent {
			//if you set ttl to -1, remove the ttl from key
			if(ttl == -1) {
				client.Persist(ctx, k)
			} else { 
				client.Expire(ctx, k, time.Duration(ttl) * time.Second)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	log.Println("Setting value(s) done!")
}
func getValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet{
		log.Println("Only GET supported for this url!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	client := redis.NewClient(&redis.Options{
                Addr:     "localhost:6379",
                Password: "", // no password set
                DB:               0,  // use default DB
        })

        ctx := context.Background()
	keyMap := r.URL.Query()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\n"))
	for _, keys := range keyMap {
		for _, key := range keys {
			//get all values associated with key
			values, err1 := client.LRange(ctx, key, 0, -1).Result()
			if err1 != nil {
				log.Println(err1.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if len(values) > 0 {
				w.Write([]byte(key + ": {\n"))
				for _, val := range values {
					w.Write([]byte(val + ",\n"))
				}
				remainingTime := client.TTL(ctx, key).Val()
				if remainingTime >= 0 {
					w.Write([]byte("TTL: " + remainingTime.String() + "\n"))
				}
				w.Write([]byte("},\n"))
			}
		}	
	}
	w.Write([]byte("}\n"))
	log.Println("Getting values done!")
}
