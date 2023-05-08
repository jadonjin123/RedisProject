package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/redis/go-redis/v9"
	"context"
	"reflect"
	"log"
)

func TestGet(t *testing.T) {
	//Test 1
	request := httptest.NewRequest("PUT", "http://localhost:3333/api/get", nil)
	responseRecorder := httptest.NewRecorder()
	getValue(responseRecorder, request)
	if(responseRecorder.Code != http.StatusBadRequest) {
		t.Errorf("Want status '%d', got '%d'", http.StatusBadRequest, responseRecorder.Code)
	}

	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
                Addr:     "localhost:6379",
                Password: "", // no password set
                DB:               0,  // use default DB
        })
	//Test 2 (insert then get)
	request = httptest.NewRequest("POST", "http://localhost:3333/api/insert?myschedule=GetUp&myschedule=Eat&myschedule=Drink&foo=bar&foo=foobar", nil)
	responseRecorder = httptest.NewRecorder()
	insertValue(responseRecorder, request) 
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Want status '%d', got '%d'", http.StatusOK, responseRecorder.Code)
	}

	request = httptest.NewRequest("GET", "http://localhost:3333/api/get?key1=myschedule&key2=foo", nil)
	responseRecorder = httptest.NewRecorder()
	getValue(responseRecorder, request) 
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Want status '%d', got '%d'", http.StatusOK, responseRecorder.Code)
	}
	expected := "{\nmyschedule: {\nGetUp,\nEat,\nDrink,\n},\nfoo: {\nbar,\nfoobar,\n},\n}\n"
	if !reflect.DeepEqual(responseRecorder.Body.Bytes(), []byte(expected)) {
		t.Errorf("Want body '%v', got '%v'", []byte(expected), responseRecorder.Body.Bytes())
	}

	//clear database
	client.FlushAll(ctx)
}

func TestDelete(t *testing.T) {	

	//Test 1
	request := httptest.NewRequest("GET", "http://localhost:3333/api/delete", nil)
	responseRecorder := httptest.NewRecorder()
	deleteValue(responseRecorder, request)
	if(responseRecorder.Code != http.StatusBadRequest) {
		t.Errorf("Want status '%d', got '%d'", http.StatusBadRequest, responseRecorder.Code)
	}

	//Test 2 (insert then delete)
	request = httptest.NewRequest("POST", "http://localhost:3333/api/insert?myschedule=GetUp&myschedule=Eat&myschedule=Drink&foo=bar&foo=foobar", nil)
	responseRecorder = httptest.NewRecorder()
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
                Addr:     "localhost:6379",
                Password: "", // no password set
                DB:               0,  // use default DB
        })
	insertValue(responseRecorder, request) 
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Want status '%d', got '%d'", http.StatusOK, responseRecorder.Code)
	}

	request = httptest.NewRequest("DELETE", "http://localhost:3333/api/delete?key1=myschedule", nil)
	responseRecorder = httptest.NewRecorder()
	deleteValue(responseRecorder, request) 
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Want status '%d', got '%d'", http.StatusOK, responseRecorder.Code)
	}
	values1, err1 := client.LRange(ctx, "myschedule", 0, -1).Result()
        if err1 != nil {
        	log.Println(err1.Error())
                return
        }
	val1 := []string{}
	if !reflect.DeepEqual(values1, val1) {
		t.Errorf("Lists 1 are not equal!")
	}

	values2, err2 := client.LRange(ctx, "foo", 0, -1).Result()
        if err2 != nil {
        	log.Println(err2.Error())
                return
        }

	val2 := []string{"bar", "foobar"}
	if !reflect.DeepEqual(values2, val2) {
		t.Errorf("Lists 2 are not equal!")
	}

	//clear database
	client.FlushAll(ctx)
}

func TestInsert(t *testing.T){

	//Test 1
	request := httptest.NewRequest("DELETE", "http://localhost:3333/api/insert", nil)
	responseRecorder := httptest.NewRecorder()
	insertValue(responseRecorder, request)
	if(responseRecorder.Code != http.StatusBadRequest) {
		t.Errorf("Want status '%d', got '%d'", http.StatusBadRequest, responseRecorder.Code)
	}

	//Test 2
	request = httptest.NewRequest("GET", "http://localhost:3333/api/insert", nil)
	responseRecorder = httptest.NewRecorder()
	insertValue(responseRecorder, request)
	if(responseRecorder.Code != http.StatusBadRequest) {
		t.Errorf("Want status '%d', got '%d'", http.StatusBadRequest, responseRecorder.Code)
	}

	//Test 3
	request = httptest.NewRequest("POST", "http://localhost:3333/api/insert?ttl=10&ttl=10", nil)
	responseRecorder = httptest.NewRecorder()
	insertValue(responseRecorder, request)
	if(responseRecorder.Code != http.StatusBadRequest) {
		t.Errorf("Want status '%d', got '%d'", http.StatusBadRequest, responseRecorder.Code)
	}

	//Test 4
	request = httptest.NewRequest("POST", "http://localhost:3333/api/insert?ttl=abc", nil)
	responseRecorder = httptest.NewRecorder()
	insertValue(responseRecorder, request)
	if(responseRecorder.Code != http.StatusBadRequest) {
		t.Errorf("Want status '%d', got '%d'", http.StatusBadRequest, responseRecorder.Code)
	}

	//Test 5
	request = httptest.NewRequest("POST", "http://localhost:3333/api/insert?myschedule=GetUp&myschedule=Eat&myschedule=Drink&foo=bar&foo=foobar", nil)
	responseRecorder = httptest.NewRecorder()
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
                Addr:     "localhost:6379",
                Password: "", // no password set
                DB:               0,  // use default DB
        })
	insertValue(responseRecorder, request) 
	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Want status '%d', got '%d'", http.StatusOK, responseRecorder.Code)
	}
	values1, err1 := client.LRange(ctx, "myschedule", 0, -1).Result()
        if err1 != nil {
        	log.Println(err1.Error())
                return
        }
	val1 := []string{"GetUp", "Eat", "Drink"}
	if !reflect.DeepEqual(values1, val1) {
		t.Errorf("Lists 1 are not equal!")
	}

	values2, err2 := client.LRange(ctx, "foo", 0, -1).Result()
        if err2 != nil {
        	log.Println(err2.Error())
                return
        }

	val2 := []string{"bar", "foobar"}
	if !reflect.DeepEqual(values2, val2) {
		t.Errorf("Lists 2 are not equal!")
	}

	//Test 6
	request = httptest.NewRequest("POST", "http://localhost:3333/api/insert?myschedule=GetUp&myschedule=Eat&myschedule=Drink&foo=bar&foo=foobar", nil)
	responseRecorder = httptest.NewRecorder()
	insertValue(responseRecorder, request)

	request = httptest.NewRequest("PUT", "http://localhost:3333/api/replace?myschedule=GetUp&myschedule=Eat&myschedule=Drink&foo=bar&foo=foobar", nil)
	responseRecorder = httptest.NewRecorder()
	insertValue(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Want status '%d', got '%d'", http.StatusOK, responseRecorder.Code)
	}
	values1, err1 = client.LRange(ctx, "myschedule", 0, -1).Result()
        if err1 != nil {
        	log.Println(err1.Error())
                return
        }
	val1 = []string{"GetUp", "Eat", "Drink"}
	if !reflect.DeepEqual(values1, val1) {
		t.Errorf("Lists 1 are not equal!")
	}

	values2, err2 = client.LRange(ctx, "foo", 0, -1).Result()
        if err2 != nil {
        	log.Println(err2.Error())
                return
        }

	val2 = []string{"bar", "foobar"}
	if !reflect.DeepEqual(values2, val2) {
		t.Errorf("Lists 2 are not equal!")
	}

	//clear database of testing values	
	client.FlushAll(ctx)
}
