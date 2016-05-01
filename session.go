package myapp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"os"
)


func newVisitor() *http.Cookie {
	mm := initModel()
	id, _ := uuid.NewV4()
	return makeCookie(mm, id.String())
}

func currentVisitor(m model, id string) *http.Cookie {
	mm := marshalModel(m)
	return makeCookie(mm,id)
}

func makeCookie(mm []byte, id string) *http.Cookie{
	b64 := base64.URLEncoding.EncodeToString(mm)
	code := getCode(b64)
	cookie := &http.Cookie{
		Name:  "session-ferret",
		Value: id + "|" + b64 + "|" + code,
		// Secure: true,
		HttpOnly: true,
	}
	return cookie
}

func marshalModel(m model) []byte{
	bs, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("error: ", err)
	}
	return bs
}

func initModel() []byte {
	m := model{
		Name: "",
		Pass: "",
		State: false,
		Pictures: []string{},
	}

	return marshalModel(m)
}

func userExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if(err == nil){
		return true, nil
	}
	if os.IsNotExist(err){
		return false, nil
	}
	return true, err
}
