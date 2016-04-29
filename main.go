package myapp


import (
	"html/template"
<<<<<<< HEAD
	"net/http"
	"encoding/json"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
=======
	"strings"
>>>>>>> origin/master
)

const gcsBucket = "csci-130group.appspot.com"

type Word struct {
	Name string
}

var tpl* template.Template

func init(){
	http.HandleFunc("/", index)
	http.HandleFunc("/api/check", wordCheck)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets/"))))

	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

<<<<<<< HEAD
func index(res http.ResponseWriter, req* http.Request){
	if req.Method == "POST" {
=======
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.Handle("/favicon.ico", http.NotFoundHandler())
>>>>>>> origin/master

		var w Word
		w.Name = req.FormValue("new-word")

		ctx := appengine.NewContext(req)
		log.Infof(ctx, "WORD SUBMITTED: %v", w.Name)

<<<<<<< HEAD
		key := datastore.NewKey(ctx, "Dictionary", w.Name, 0, nil)
		_, err := datastore.Put(ctx, key, &w)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	tpl.ExecuteTemplate(res, "index.html", nil)
}


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
}
=======
func index(res http.ResponseWriter, req* http.Request){
	//if req.URL.Path != "/" {
	//	http.NotFound(res, req)
	//	return
	//}
	cookie := genCookie(res, req)
	http.SetCookie(res, cookie)
	tpl.ExecuteTemplate(res, "index.html", nil)
}

func login(res http.ResponseWriter, req *http.Request){
	cookie := genCookie(res, req)
	http.SetCookie(res, cookie)
	tpl.ExecuteTemplate(res, "login.html", nil)
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
>>>>>>> origin/master
