package myapp


import (
	"html/template"
	"strings"
	//"io"
	"net/http"
	"fmt"
	"google.golang.org/cloud/storage"
	//"encoding/json"
	"google.golang.org/appengine/urlfetch"

	"encoding/json"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"io"
	"io/ioutil"
)

//we'll use this bucket constant for uploading files to our bucket
const gcsBucket = "user-images"

var tpl* template.Template

func init(){

	//first we parse our html and serve our css files.
	//since matt loves local pics, the pictures are also being served....
	// but we are using gcs links from our ferret bucket.
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets/"))))



	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/register", register)
	http.HandleFunc("/changepass", changePass)
	http.HandleFunc("/files", files)
	http.HandleFunc("/photos", gallery)
	http.HandleFunc("/api/check", wordCheck)
	http.HandleFunc("/gif", gif)
	http.HandleFunc("/giffy", giffy)
	http.Handle("/favicon.ico", http.NotFoundHandler())


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
		cookie = getPhotos(cookie,req,res)
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
		if(user.Name == req.FormValue("name") || req.FormValue("name") == "" || req.FormValue("password") == ""){
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

func changePass(res http.ResponseWriter, req *http.Request){
	cookie := genCookie(res, req)
	if req.Method == "POST"{
		m := Model(cookie)
		if (m.Pass == "") ||
		   (req.FormValue("password") != m.Pass) ||
		   (req.FormValue("password2") != req.FormValue("password3")){
			fmt.Sprintf("Wrong!")
			http.Redirect(res, req, "/changepass", 302)
			return
		}
		m.Pass = req.FormValue("password2")
		m.State = false
		setUser(req, m)
		m.State = true
		xs := strings.Split(cookie.Value, "|")
		id := xs[0]

		cookie := currentVisitor(m, id)
		http.SetCookie(res, cookie)
		http.Redirect(res, req, "/", 302)
		return
	}
	tpl.ExecuteTemplate(res, "changepass.html", nil)
}

func files(res http.ResponseWriter, req *http.Request){
	cookie := genCookie(res, req)
	if req.Method == "POST" {
		src, hdr, err := req.FormFile("data")
		if err != nil{
			panic(err)
		}
		cookie = uploadPhoto(src, hdr, cookie, req)
		http.SetCookie(res, cookie)
	}
	m := Model(cookie)
	tpl.ExecuteTemplate(res, "files.html", m)
}

func gallery(res http.ResponseWriter, req *http.Request){
	cookie := genCookie(res, req)
	m := Model(cookie)
	tpl.ExecuteTemplate(res, "gallery.html", m)
}

func gif(res http.ResponseWriter, req *http.Request){
	cookie := genCookie(res, req)
	m := Model(cookie)
	if req.Method == "POST"{
		http.Redirect(res, req, "/giffy", 302)
		return
	}
	tpl.ExecuteTemplate(res, "gif.html", m)
}

func wordCheck(res http.ResponseWriter, req *http.Request) {

	ctx := appengine.NewContext(req)

	// acquire the incoming word
	var m model
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Infof(ctx, err.Error())
	}
	name := string(bs)
	log.Infof(ctx, "ENTERED wordCheck - name: %v", name)

	// check the incoming word against the datastore
	key := datastore.NewKey(ctx, "Users", name, 0, nil)
	err = datastore.Get(ctx, key, &m)
	//m = getUser(req, name)
	if m.Name == name {
		io.WriteString(res, "true")
		return
	}
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

//Load the images that belong in the folder with the users name
//I ADDED A PARAMETER THE REQUEST and RESPONSE
func getPhotos(c *http.Cookie, req *http.Request, res http.ResponseWriter) *http.Cookie{
	//Get the model from the cookie to prepare to append pictures
	m := Model(c)

	//grabbing context for error checking
	ctx := appengine.NewContext(req)

	//creating a new client from our context
	client, err := storage.NewClient(ctx)
	if err != nil{
		log.Errorf(ctx, "Error in getPhotos client err")
	}
	defer client.Close()

	//grabbing a client fronm our specific bucket
	bucket := client.Bucket(gcsBucket)

	//making a query so we find a specific user's files
	q := &storage.Query{
		Prefix: m.Name,
	}

	//Store all filenames that are in the user's folder in a slice of strings
	listObjects, err := bucket.List(ctx, q)
	if err != nil {
		log.Errorf(ctx, "getPhotos: unable to list bucket %q: %v", gcsBucket, err)
		return c
	}
	//for loop to append all file paths
	for _, obj := range listObjects.Results {

		ext := obj.Name[strings.LastIndex(obj.Name, ".")+1:]
		//log.Infof(d.ctx, "%v", ext)
		if ext == "jpg" || ext == "jpeg" || ext == "png"{
			m.Pictures = append(m.Pictures, obj.MediaLink)
		}
		m.Files = append(m.Files, obj.MediaLink)
	}

	xs := strings.Split(c.Value, "|")
	id := xs[0]
	
	c = currentVisitor(m, id)

	return c
}


func giffy(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	t := req.FormValue("term") + "_ferret"

	client := urlfetch.Client(ctx)
	result, err := client.Get("http://api.giphy.com/v1/gifs/search?q=" + t + "&api_key=dc6zaTOxFJmzC")
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	defer result.Body.Close()
	var obj struct {
		Data []struct {
			URL    string `json:"url"`
			Images struct {
					   Original struct {
									URL string
								}
				   }
		}
	}
	err = json.NewDecoder(result.Body).Decode(&obj)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	for _, img := range obj.Data {
		fmt.Fprintf(res, `<a href="%v"><img src="%v"></a><br>`, img.URL, img.Images.Original.URL)
	}
	cookie := genCookie(res, req)
	m := Model(cookie)
	tpl.ExecuteTemplate(res,"gifs.html", m)
}
