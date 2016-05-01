package myapp

import(
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"net/http"
	"log"
)

type model struct {
	Name string
	Pass string
	State bool
	Pictures []string
}

func Model(c *http.Cookie) model {
	xs := strings.Split(c.Value, "|")
	usrData := xs[1]

	bs, err := base64.URLEncoding.DecodeString(usrData)
	if err != nil{
		log.Println("Error decoding base64", err)
	}

	var m model
	err = json.Unmarshal(bs, &m)
	if err != nil{
		fmt.Printf("error unmarshalling: %v", err)
	}
	return m
}


func AltModel(usrData string) model {
	bs, err := base64.URLEncoding.DecodeString(usrData)
	if err != nil{
		log.Println("Error decoding base64", err)
	}

	var m model
	err = json.Unmarshal(bs, &m)
	if err != nil{
		fmt.Printf("error unmarshalling: %v", err)
	}
	return m
}
