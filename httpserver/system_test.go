package main

import (
	"testing"
	"log"
	"os/exec"
	"strings"
)

func TestNotFound(t *testing.T) {

	//Test 1
	curl := exec.Command("curl", "-v", "http://localhost:3333/api/abcd")
	out, err := curl.Output()
	if err != nil {
		log.Println("ERROR")
		log.Println(err.Error())
	}
	if !strings.Contains(string(out), "404") {
		t.Errorf("Error code is not 404!")
	}

	//Test 2
	curl = exec.Command("curl", "-v", "-X", "DELETE", "http://localhost:3333/api/abcd")
	out, err = curl.Output()
	if err != nil {
		log.Println("ERROR")
		log.Println(err.Error())
	}
	if !strings.Contains(string(out), "404") {
		t.Errorf("Error code is not 404!")
	}

	//Test 3
	curl = exec.Command("curl", "-v", "-X", "POST", "http://localhost:3333/api/abcd")
	out, err = curl.Output()
	if err != nil {
		log.Println("ERROR")
		log.Println(err.Error())
	}
	if !strings.Contains(string(out), "404") {
		t.Errorf("Error code is not 404!")
	}

	//Test 4
	curl = exec.Command("curl", "-v", "-X", "PUT", "http://localhost:3333/api/abcd")
	out, err = curl.Output()
	if err != nil {
		log.Println("ERROR")
		log.Println(err.Error())
	}
	if !strings.Contains(string(out), "404") {
		t.Errorf("Error code is not 404!")
	}
}
