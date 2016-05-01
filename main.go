package myapp


import (
	"html/template"
	"strings"
	"io"
	"net/http"
	//"encoding/json"
	//"google.golang.org/appengine"
	//"google.golang.org/appengine/datastore"
	//"google.golang.org/appengine/log"
)

const gcsBucket = "csci-130group.appspot.com"

type Word struct {
	Name string
}

var tpl* template.Template

func init(){
	//http.HandleFunc("/", index)
	//http.HandleFunc("/api/check", wordCheck)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets/"))))

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.Handle("/favicon.ico", http.NotFoundHandler())
/*func index(res http.ResponseWriter, req* http.Request){
	if req.Method == "POST" {
>>>>>>> 70880bd3b3e1aa5d9fa834803870bc6e25c90ac2

		var w Word
		w.Name = req.FormValue("new-word")

		ctx := appengine.NewContext(req)
		log.Infof(ctx, "WORD SUBMITTED: %v", w.Name)

*/
}
func index(res http.ResponseWriter, req* http.Request){
	//if req.URL.Path != "/" {
	//	http.NotFound(res, req)
	//	return
	//}
	cookie := genCookie(res, req)
	m := Model(cookie)
	m.State = true
	xs := strings.Split(cookie.Value, "|")
	id := xs[0]
	cookie = currentVisitor(m, id)
	http.SetCookie(res, cookie)
/*
		key := datastore.NewKey(ctx, "Dictionary", w.Name, 0, nil)
		_, err := datastore.Put(ctx, key, &w)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
>>>>>>> 70880bd3b3e1aa5d9fa834803870bc6e25c90ac2*/
	tpl.ExecuteTemplate(res, "index.html", nil)
	io.WriteString(res, "UUID: " + id)
}

func login(res http.ResponseWriter, req *http.Request){
	cookie := genCookie(res, req)
	/*if req.Method == "POST"{
		exists, _ := userExists(req.FormValue("name"))
		if(exists == false){
			http.Redirect(res, req, "/login", 302)
			return
		}
		mod := getUser(req.FormValue("name"))
		if mod.Pass != req.FormValue("password"){
			http.Redirect(res, req, "/login", 302)
			return
		}
		xs := strings.Split(cookie.Value, "|")
		id := xs[0]
		mod.State = true
		cookie = currentVisitor(mod, id)
		http.SetCookie(res, cookie)

	}*/
	m := Model(cookie)
	xs := strings.Split(cookie.Value, "|")
	id := xs[0]
	http.SetCookie(res, cookie)
	tpl.ExecuteTemplate(res, "login.html", nil)
	if(m.State == true){
		io.WriteString(res, "UUID: " + id)
	}
}

func register(res http.ResponseWriter, req *http.Request){
	cookie := genCookie(res, req)
	/*if req.Method == "POST"{
		exists, _ := userExists(req.FormValue("name"))
		if(exists){
			http.Redirect(res, req, "/register", 302)
			return
		}
		m := Model(cookie)
		m.State = false;
		m.Name = req.FormValue("name")
		m.Pass = req.FormValue("password")
		xs := strings.Split(cookie.Value, "|")
		id := xs[0]
		writeFile(cookie);
		m.State = true;
		cookie := currentVisitor(m, id)
		http.SetCookie(res, cookie)

		http.Redirect(res, req, "/", 302)
		return
	}*/
	m := Model(cookie)
	xs := strings.Split(cookie.Value, "|")
	id := xs[0]
	http.SetCookie(res, cookie)
	tpl.ExecuteTemplate(res, "register.html", nil)
	if(m.State == true){
		io.WriteString(res, "UUID: " + id)
	}
}


func genCookie(res http.ResponseWriter, req *http.Request) *http.Cookie{
	cookie, err := req.Cookie("session-fino")
	if err != nil{
		cookie = newVisitor()
		http.SetCookie(res, cookie)
	}
	if strings.Count(cookie.Value, "|") != 2{
		cookie = newVisitor()
		http.SetCookie(res, cookie)
	}
	if tampered(cookie.Value){
		cookie = newVisitor()
		http.SetCookie(res, cookie)
	}
	return cookie
}



/*>>>>>>> 70880bd3b3e1aa5d9fa834803870bc6e25c90ac2

func wordCheck(res http.ResponseWriter, req *http.Request) {

	ctx := appengine.NewContext(req)

	// acquire the incoming word
	var w Word
	json.NewDecoder(req.Body).Decode(&w)
	log.Infof(ctx, "ENTERED wordCheck - w.Name: %v", w.Name)

	// check the incoming word against the datastore
	key := datastore.NewKey(ctx, "Dictionary", w.Name, 0, nil)
	err := datastore.Get(ctx, key, &w)
	if err != nil {
		json.NewEncoder(res).Encode("false")
		return
	}
	json.NewEncoder(res).Encode("true")
}*/
