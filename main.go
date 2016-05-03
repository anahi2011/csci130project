package myapp


import (
	"html/template"
	"strings"
	//"io"
	"net/http"
	"fmt"
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

	//http.HandleFunc("/api/check", wordCheck)

	//first we parse our html and serve our css files.
	//since matt loves local pics, the pictures are also being served....
	// but we are using gcs links from our ferret bucket.
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets/"))))



	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
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

//every time the user loads the page they get a new uuid....
//they shouldn't get a new one...
func index(res http.ResponseWriter, req* http.Request){
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	//so when user enters main webpage we make a cookie. ok
	cookie := genCookie(res, req)


	m := Model(cookie)

	//we split up the values in our cookie by the delimiter |

	//remember our cookie value is set up like this
	// uuid | modelEncodeToB64 | HMAC



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
	tpl.ExecuteTemplate(res, "index.html", m)
}

//wtf is going on with this login....
func login(res http.ResponseWriter, req *http.Request){
	cookie := genCookie(res, req)
	if req.Method == "POST"{
		if(req.FormValue("name") == "" || req.FormValue("password") == ""){
			fmt.Printf("Type in something dipshit")
			http.Redirect(res, req, "/login", 302)
			return
		}
		m := getUser(req, req.FormValue("name"))
		if m.Pass == ""{
			fmt.Printf("User Not Found")
			http.Redirect(res, req, "/login", 302)
			return
		}
		if(m.Pass != req.FormValue("password")){
			fmt.Printf("Wrong Password Dumbass")
			http.Redirect(res, req, "/login", 302)
			return
		}
		m2 := Model(cookie)
		m2.State = true
		m2.Name = m.Name
		m2.Pass = m.Pass
		m2.Pictures = m.Pictures
		xs := strings.Split(cookie.Value, "|")
		id := xs[0]

		cookie := currentVisitor(m2, id)
		http.SetCookie(res, cookie)

		http.Redirect(res, req, "/", 302)
		return
	}
	tpl.ExecuteTemplate(res, "login.html", nil)
}

func logout(res http.ResponseWriter, req *http.Request){
	cookie := newVisitor()
	http.SetCookie(res, cookie)
	http.Redirect(res, req, "/", 302)
}

//:( the lack of comments makes me want to cry
func register(res http.ResponseWriter, req *http.Request){
	cookie := genCookie(res, req)
	if req.Method == "POST"{
		user := getUser(req, req.FormValue("name"))
		if(user.Name == req.FormValue("name")){
			fmt.Printf("Username already taken")
			http.Redirect(res, req, "/register", 302)
			return
		}
		m := Model(cookie)
		m.Name = req.FormValue("name")
		m.Pass = req.FormValue("password")
		setUser(req, m)
		m.State = true
		xs := strings.Split(cookie.Value, "|")
		id := xs[0]

		cookie := currentVisitor(m, id)
		http.SetCookie(res, cookie)
		http.Redirect(res, req, "/", 302)
		return
	}
	tpl.ExecuteTemplate(res, "register.html", nil)
}


//looks for cookie and returns it.
//if it doesn't exits we make a new one and then sets it and returns it.
func genCookie(res http.ResponseWriter, req *http.Request) *http.Cookie{

	cookie, err := req.Cookie("session-ferret")
	if err != nil{
		cookie = newVisitor()
		http.SetCookie(res, cookie)
		//return cause if we made the cookie... welll theres no need to
		//check if it was tampered...
		return cookie
	}

	return cookie

	if strings.Count(cookie.Value, "|") != 2{
		cookie = newVisitor()
		http.SetCookie(res, cookie)
	}

	//if the user fucked up the cookie we make a new one
	//we test the cookie using the HMAC code
	//if you don't know how hmac works. well..
	//it's not hard. we generate a hmac code with our function.
	//we put it in our cookie with the data that got hmaced. so when
	//we get the cookie back we hmac the data and compare it to the hmac code.
	//the reason this works for us to verify is because the user doesn't know the secret
	//key we use to make our hmac code. Since matt chose a super secure key, nobody should
	//be able to crack it.
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
