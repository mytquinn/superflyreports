package main

import
(
	"net/http"
	"io/ioutil"
	"strings"
	"log"
    "userAuth"
	"formHandler"
	)

func main() {
	myMux := http.NewServeMux()
	//Setup Routing
	myMux.HandleFunc("/", ServeHTTP)
	myMux.HandleFunc("/doSignout", userAuth.SignOut)
	myMux.HandleFunc("/doSignin", userAuth.DoSignin)
	myMux.HandleFunc("/doSignup", userAuth.DoSignup)
	myMux.HandleFunc("/doRecoverPassword", userAuth.DoRecoverPassword)
	myMux.HandleFunc("/verifyEmail/", userAuth.VerifyEmail)
	myMux.HandleFunc("/checkLogin", userAuth.CheckLogin)
	myMux.HandleFunc("/getSessionData/", userAuth.GetSessionData)
	myMux.HandleFunc("/getType/", formHandler.GetType)
	myMux.HandleFunc("/journalSubmit", formHandler.JournalSubmit)
	//Start server
	err := http.ListenAndServe(":80", myMux)
    if err != nil {
		log.Println(err.Error())
	}
}


func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path[1:]

	directory := strings.Split(path, "/")[0]

	//Ignore path and serve index.html unless it starts with these directories
	if directory != "" {
		indexList := map[string] bool {
			"views": true,
			"images" : true,
			"css" : true,
			"dart" : true,
			"fonts" : true,
		}

		if !indexList[directory]  {
			path = "index.html"
		}

	} else {
		path = "index.html"
	}


	data, err := ioutil.ReadFile(string(path))

    // Set MINE types for file transfer
	if err == nil {
		var contentType string

		if strings.HasSuffix(path, ".css") {
			contentType = "text/css"
		} else if strings.HasSuffix(path, ".html") {
			contentType = "text/html"
		} else if strings.HasSuffix(path, ".jpg") {
			contentType = "image/jpeg"
		} else if strings.HasSuffix(path, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(path, ".js") {
			contentType = "application/javascript"
		} else if strings.HasSuffix(path, ".dart") {
			contentType = "application/dart"
		} else if strings.HasSuffix(path, ".gif") {
			contentType = "text/gif"
		} else if strings.HasSuffix(path, ".eot") {
			contentType = "application/vnd.ms-fontobject"
		} else if strings.HasSuffix(path, ".svg") {
			contentType = "image/svg+xml"
		}else if strings.HasSuffix(path, ".ttf") {
			contentType = "application/x-font-ttf"
		}else if strings.HasSuffix(path, ".woff") {
			contentType = "application/x-font-woff"
		} else if strings.HasSuffix(path, ".woff2") {
			contentType = "application/x-font-woff"
		} else {
			contentType = "text/plain"
		}
		w.Header().Add("Content-Type", contentType)
		w.Write(data)
	} else {
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
	}
}