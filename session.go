package myapp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"os"
)

//we make an object for the new user
//then it's marshalled to json string
//we assign a uuid to this new user
//then we make a cookie that has
// uuid|modelEncodeToB64|HMAC as its value
func newVisitor() *http.Cookie {

	//mm is a string. it is a json string of the model data
	mm := initModel()
	id, _ := uuid.NewV4()

	return makeCookie(mm, id.String())
}

//returns a cookie that is made with the passed model and uuid
func currentVisitor(m model, id string) *http.Cookie {
	mm := marshalModel(m)
	return makeCookie(mm,id)
}

//encodes cookie to base 64 cause we can't have multiple (of some symbol)
// "" i think the symbol is the "" or
//cause golang is gay and freaks out or something
//then it is putting the uuid, the b64 encoded model, and the HMAC code.
//ok. these 3 things are put together split up by the | pipe.
//i think they call it a delimiter.. not sure though. then we assign
//the 3 things together to the cookies value
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

//takes a model and marshals it to jason
//returns the string of jason data i guess
func marshalModel(m model) []byte{
	bs, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("error: ", err)
	}
	return bs
}

//makes a defaulted model struct instance(aka object)
//returns a string of the object marshalled to json
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
