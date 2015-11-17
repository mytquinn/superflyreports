package userAuth

import
(
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"log"
	"github.com/astaxie/beego/session"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sendgrid/sendgrid-go"
	"time"
	"strings"
	"config"
)

var globalSessions *session.Manager
var anonSessions *session.Manager
var db *sql.DB

//Handles sign functionality from form POST or page GET
//Sets user_id and userSession_id cookies and saves username and user_id to the session on successful login
func DoSignin(w http.ResponseWriter, r *http.Request) {

	var id int


	if r.Method == "POST" {
		// Handles login when it is hit as a post request
		r.ParseForm()
		stmt, err := db.Prepare("select id from users where username=? and password=?")
		res := stmt.QueryRow(r.FormValue("username"), r.FormValue("password"))
		err = res.Scan(&id)

		if (err == nil) {
			sess, _ := globalSessions.SessionStart(w, r)
			defer sess.SessionRelease(w)
			setUserCookies(w, id, sess.SessionID())
			_ = sess.Set("user_id", id)
			_ = sess.Set("username", r.FormValue("username"))
			if r.FormValue("remember-me") == "on" {
				saveSession(w, r, sess.SessionID(), id)

			}
			addRemoteAddress(r, id)
			http.Redirect(w, r, "/", 302)
		} else {
			log.Println("Database connection failed: ", err)
		}
	} else {
		anonsess, _ := anonSessions.SessionStart(w, r)
		defer anonsess.SessionRelease(w)
		// Handles auto login when it is hit as a GET request
		sessionIdCookie, err := r.Cookie("userSession_id")
		if err == nil {
			stmt, err := db.Prepare("select id, username from users where session_id=?")
			res:= stmt.QueryRow(sessionIdCookie.Value)
			var username string
			err = res.Scan(&id, &username)
			if err == nil{
				if checkRemoteAddress(r, id) {
					sess, _ := globalSessions.SessionStart(w, r)
					defer sess.SessionRelease(w)
					err = sess.Set("user_id", id)
					if err != nil {
						log.Println(err)
					}
					_ = sess.Set("username", username)
					saveSession(w, r, sess.SessionID(), id)
					setUserCookies(w, id, sess.SessionID())
					http.Redirect(w, r, "/", 302)
				} else {
					http.Redirect(w, r, "/newAddress", 302)
				}
			} else {
				http.Redirect(w, r, "/userNotFound", 302)
			}
		} else {
			http.Redirect(w, r, "/", 302)
		}
	}
}

// Add the users ip address to database to be used during auto login
func addRemoteAddress(r *http.Request, user_id int) {
	var addresses string
	stmt, _ := db.Prepare("select remote_addr from users where id=?")
	res:= stmt.QueryRow(user_id)
	err := res.Scan(&addresses)
	newAddr := strings.Split(r.RemoteAddr, ":")[0]
	if newAddr == "[" {
		newAddr = "localHost"
	}
	if !strings.Contains(addresses, newAddr) && (err == nil || addresses == "") {
		if addresses != "" {
			addresses = addresses + ";" + newAddr
		} else {
			addresses = newAddr
		}

		stmt, _ := db.Prepare("update users set remote_addr=? where id=?")
		stmt.Exec(addresses, user_id)
	}

}

// Checks the database to makes sure the user has logged in from address before.
func checkRemoteAddress(r *http.Request, user_id int) bool {
	var addresses string
	stmt, _ := db.Prepare("select remote_addr from users where id=?")
	res:= stmt.QueryRow(user_id)
	err := res.Scan(&addresses)
	currentAddr := strings.Split(r.RemoteAddr, ":")[0]
	if currentAddr == "[" {
		currentAddr = "localHost"
	}
	if strings.Contains(addresses, currentAddr) && err == nil {
		return true
	} else {
		return false
	}
}

// Writes session to user record for auto login
func saveSession(w http.ResponseWriter, r *http.Request, sid string, user_id int){

	stmt, err := db.Prepare("update users set session_id=? where id=?")
	_, err = stmt.Exec(sid, user_id)
	if err != nil  {
		log.Println("Update session_id failed: ", err)
	}
}

// Sets cookies with user_id and session id to be used for auto login
func setUserCookies(w http.ResponseWriter, id int, sessId string){
	userIdCookie := http.Cookie{
		Name : "user_id",
		Value : strconv.Itoa(id),
		Expires : time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: false,
		Path: "/",
	}
	userSessionCookie := http.Cookie{
		Name : "userSession_id",
		Value : sessId,
		Expires : time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: false,
		Path: "/",
	}
	http.SetCookie(w, &userIdCookie)
	http.SetCookie(w, &userSessionCookie)
}

// Processes for data for signup and sends email to verify account
func DoSignup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vKey := make([]byte, 32)
	n, err := rand.Read(vKey)
	if n != len(vKey) || err != nil {
		log.Println("Could not successfully read from the system CSPRNG.")
	}
	validationKey := hex.EncodeToString(vKey)
	stmt, _ := db.Prepare("insert into signup(username, email, password, validationKey) values(?,?,?,?)")
	_, err = stmt.Exec(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"), validationKey)
	if err != nil {
		// if a validation requests already exists resend email
		if strings.Contains(err.Error(), "1062") {
			log.Println("1062 error")
			stmt, _ := db.Prepare("select validationKey from signup where username=?")
			res := stmt.QueryRow(r.FormValue("username"))
			res.Scan(&validationKey)
			sendVerification(r.FormValue("email"), validationKey)
			http.Redirect(w, r, r.URL.Host + "/resendValidation", 302)
		} else {
			log.Print("Error creating signup record")
			log.Println(err)
		}
	} else {
		sendVerification(r.FormValue("email"), validationKey)
		http.Redirect(w, r, r.URL.Host + "/validationSent", 302)
	}
}

// Sends verification email to the user when they signup
func sendVerification(email string, validationKey string) {
	sg := sendgrid.NewSendGridClient(myConfig.Sg_user, myConfig.Sg_password)
	message := sendgrid.NewMail()
	message.AddTo(email)
	message.SetSubject("Verify your new Super Fly Reports Account")
	message.SetText("Please follow or copy this link to verify your new account:" + myConfig.Sg_emailLink + validationKey)
	message.SetHTML(`<html><body>Please follow or copy this link to verify your new account:
                      <a href='` + myConfig.Sg_emailLink + validationKey +
	`'>`+ myConfig.Sg_emailLink + validationKey + `</a></body></html>`)
	message.SetFrom("account_verify@superfly.com")
	if err := sg.Send(message); err == nil {
		log.Println("Email sent to :", email)
	} else {
		log.Println(err)
	}
}

// Handles verification of account email based on validationKey
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	validationKey := strings.Split(r.RequestURI, "/")[2]
	stmt, _ := db.Prepare("select username, email, password from signup where validationKey = ?")
	var username string
	var email string
	var password string
	res := stmt.QueryRow(validationKey)
	err := res.Scan(&username, &email, &password)
	if err != nil {
		http.Redirect(w, r, "/validationFailed", 302)
	}

	stmt, _ = db.Prepare("insert into users (username, email, password) values (?,?,?)")
	row, err := stmt.Exec(username, email, password)
	if err == nil {
		id64, _ := row.LastInsertId()
		id := int(id64)
		// Login user and delete signup record
		var sess session.SessionStore
		sessionCookie, err := r.Cookie("session_id")
		if err == nil {
			sess, _ = globalSessions.GetSessionStore(sessionCookie.Value)
		} else {
			sess, _ = globalSessions.SessionStart(w, r)
		}
		defer sess.SessionRelease(w)

		_ = sess.Set("user_id", id)
		_ = sess.Set("username", username)
		setUserCookies(w, id, sess.SessionID())
		saveSession(w, r, sess.SessionID(), id)
		addRemoteAddress(r, id)
		db.Prepare("delete from signup where validationKey = ?")
		db.Exec(validationKey);
		http.Redirect(w, r, "/", 302)
	}
}

func CheckLogin (w http.ResponseWriter, r *http.Request) {
	sessCookie, err := r.Cookie("session_id")
	if err == nil {
		sess, err := globalSessions.GetSessionStore(sessCookie.Value)
		user_id := sess.Get("user_id")
		if err == nil && user_id != nil {
			if checkRemoteAddress(r, user_id.(int)) {
				w.WriteHeader(302)
				_, _ = w.Write([]byte("true"))
				return
			}
		}
	}
	w.WriteHeader(302)
	_, _ = w.Write([]byte("false"))
}

func GetSessionData (w http.ResponseWriter, r *http.Request) {
	data := strings.Split(r.RequestURI, "/")[2]
	sessCookie, err := r.Cookie("session_id")
	if sessCookie.Value != "" && err == nil {
		sess, err := globalSessions.GetSessionStore(sessCookie.Value)
		sessData := sess.Get(data)
		if err == nil {
			w.WriteHeader(203)
			_, _ = w.Write([]byte(sessData.(string)))
		} else {
			w.WriteHeader(302)
			_, _ = w.Write([]byte("null"))
		}
	}
}

func SignOut (w http.ResponseWriter, r *http.Request) {
	globalSessions.SessionDestroy(w, r)
}

// Handles password recovery process.
func DoRecoverPassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

}

func init() {
	db, _ = sql.Open("mysql", myConfig.Db_user + ":" + myConfig.Db_password + "@" + myConfig.Db_address + "/" + myConfig.Db_schema)
	err := db.Ping()
	if err == nil {
		log.Println("DB responded")
	} else {
		log.Println("DB not responding: ", err)
	}

	dbKeepalive := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range dbKeepalive.C {
			err := db.Ping()
			if err != nil {
				log.Println("DB Connection droppped")
			}
		}
	}()

	anonSessions, err = session.NewManager("memory", `{"cookieName":"anonsession_id","gclifetime":3600}`)
	if err != nil {
		log.Println(err)
	}
	go anonSessions.GC()

	globalSessions, err = session.NewManager("memory", `{"cookieName":"session_id","gclifetime":3600}`)
	if err != nil {
		log.Println(err)
	}
	go globalSessions.GC()
}

