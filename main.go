package myapp

import (
	//"golang.org/x/net/context"
	"net/http"
	//"google.golang.org/cloud/storage"
	//"io"
	"html/template"
	"strings"
	"io"
)

const gcsBucket = "csci-130group.appspot.com"


//we could use this to pass it to our functions
//that way we don't have to have too many parameters
//NOT SURE IF WE WILL USE THIS THOUGH
//type demo struct{
//	ctx	context.Context
//	w	http.ResponseWriter
//	bucket	*storage.BucketHandle
//	client 	*storage.Client
//}

var tpl* template.Template

func init(){
	tpl = template.Must(tpl.ParseGlob("templates/*.html"))
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets/"))))


	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.Handle("/favicon.ico", http.NotFoundHandler())


	http.ListenAndServe(":8080",nil)
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
	cookie, err := req.Cookie("session-ferret")
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

//func putFile(ctx context.Context, name string, rdr io.Reader) error {
//
//	client, err := storage.NewClient(ctx)
//	if err != nil {
//		return err
//	}
//	defer client.Close()
//
//	writer := client.Bucket(gcsBucket).Object(name).NewWriter(ctx)
//
//	io.Copy(writer, rdr)
//	// check for errors on io.Copy in production code!
//	return writer.Close()
//}
//
//func getFile(ctx context.Context, name string) (io.ReadCloser, error) {
//	client, err := storage.NewClient(ctx)
//	if err != nil {
//		return nil, err
//	}
//	defer client.Close()
//
//	return client.Bucket(gcsBucket).Object(name).NewReader(ctx)
//}
//
//func getFileLink(ctx context.Context, name string) (string, error) {
//	client, err := storage.NewClient(ctx)
//	if err != nil {
//		return "", err
//	}
//	defer client.Close()
//
//	attrs, err := client.Bucket(gcsBucket).Object(name).Attrs(ctx)
//	if err != nil {
//		return "", err
//	}
//	return attrs.MediaLink, nil
//}
