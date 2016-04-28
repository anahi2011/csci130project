package myapp

import (
	//"golang.org/x/net/context"
	"net/http"
	//"google.golang.org/cloud/storage"
	//"io"
	"html/template"
	"strings"
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
	http.Handle("/favicon.ico", http.NotFoundHandler())


	http.ListenAndServe(":8080",nil)
}

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
