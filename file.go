package myapp

import (
	"net/http"
	"os"
	"bufio"
)

func writeFile(c *http.Cookie){
	f, err := os.OpenFile("users.txt", os.O_APPEND | os.O_WRONLY, 0666)
	if(err != nil){
		panic(err)
	}
	m := Model(c)
	bs := marshalModel(m)
	_, err = f.WriteString(string(bs))
	if err != nil {
		panic(err)
	}
	f.Close();
}

func searchFile(name string) bool{
	f, err := os.Open("users.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m := AltModel(scanner.Text())
		if(m.Name == name){
			return true
		}
	}
	return false;
}

func getUser(name string) model{
	f, err := os.Open("users.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		m := AltModel(scanner.Text())
		if(m.Name == name){
			return m
		}
	}
	return model{};
}
