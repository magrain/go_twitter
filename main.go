package main 
import(
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "github.com/gorilla/context"
   "net/http"
   "html/template"
   "path"
   "os"
   "strconv"
   "fmt"
   "crypto/rand"
   "encoding/hex"
   "encoding/base64"
   "io"
   "golang.org/x/crypto/scrypt"
   "database/sql"
   _ "github.com/go-sql-driver/mysql"
)

type User struct {
   id int
   name string
   display_name string
   email string
   password string
   salt string
}

var store *sessions.CookieStore
var db    *sql.DB

func getTemplatePath(file string) string {
   return path.Join("contents", file)
}

func getSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
   session, _ := store.Get(r, "gotwitter")
   return session
}

func authenticate(w http.ResponseWriter, r *http.Request, email, passwd string) {
   //query := `SELECT id, name, display_name, email FROM users where email=? AND password=?` 
   //row := db.QueryRow(`SELECT id, name, display_name, email FROM users where email=? AND password=?`, email, passwd)
   row := db.QueryRow(`SELECT id, name, display_name, email, password, salt FROM users where email=?`, email)
   user := User{}
   //err := row.Scan(&user.id, &user.name, &user.display_name, &user.email)
   err := row.Scan(&user.id, &user.name, &user.display_name, &user.email, &user.password, &user.salt)
   if err != nil {
      session := getSession(w, r)
      delete(session.Values, "user_id")
      session.Save(r, w)
      render(w, r, http.StatusUnauthorized, "login.html", nil)
      return 
   }
   // (password string, salt string, hash string) 
   if !checkHashFromPassword(passwd, user.salt, user.password) {
      session := getSession(w, r)
      delete(session.Values, "user_id")
      session.Save(r, w)
      render(w, r, http.StatusUnauthorized, "login.html", nil)
      return
   }
   session := getSession(w, r)
   session.Values["user_id"] = user.id
   session.Save(r, w)
}

func authenticated(w http.ResponseWriter, r *http.Request) bool {
   user := getCurrentUser(w, r)
   if user == nil {
      http.Redirect(w, r, "/login", http.StatusFound)
      return false
   }
   return true
}

func getCurrentUser(w http.ResponseWriter, r *http.Request) *User {
   u := context.Get(r, "user")
   if u != nil {	
      user := u.(User)
         return &user
   }
   session := getSession(w, r)
   userID, ok := session.Values["user_id"]
   if !ok || userID == nil {
      return nil
   }
   row := db.QueryRow(`SELECT id, name, display_name, email FROM users WHERE id=?`, userID)
   user := User{}
   err := row.Scan(&user.id, &user.name, &user.display_name, &user.email)
   if err == sql.ErrNoRows {
      session := getSession(w, r)
      delete(session.Values, "user_id")
      session.Save(r, w)
      render(w, r, http.StatusUnauthorized, "login.html", nil)
      return nil
   }
   context.Set(r, "user", user)
   return &user
}
func render(w http.ResponseWriter, r *http.Request, status int, file string, data interface{}){
   t := template.Must(template.New(file).ParseFiles(getTemplatePath(file)))
   w.WriteHeader(status)
   t.Execute(w, data)
}

func GetLogin(w http.ResponseWriter, r *http.Request){
   render(w, r, http.StatusOK, "login.html", nil)
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
   email := r.FormValue("email")
   passwd := r.FormValue("password")
   fmt.Println(email)
   authenticate(w, r, email, passwd)
   http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetLogout(w http.ResponseWriter, r *http.Request) {
   session := getSession(w, r)
   delete(session.Values, "user_id")
   session.Options = &sessions.Options{MaxAge: -1}
   session.Save(r, w)
   http.Redirect(w, r, "/login", http.StatusFound)
}

func GetRegister(w http.ResponseWriter, r *http.Request) {
   render(w, r, http.StatusOK, "register.html", nil)
}

func PostRegister(w http.ResponseWriter, r *http.Request){
   name := r.FormValue("name")
   display_name := r.FormValue("display_name")
   email := r.FormValue("email")
   password := r.FormValue("password")
   if name == "" || password == "" || email == "" || display_name == "" {
      http.Error(w, http.StatusText(http.StatusBadRequest), 400)
      return
   }
   userID := register(name, display_name, email ,password)
   session := getSession(w, r)
   session.Values["user_id"] = userID
   session.Save(r, w)
   http.Redirect(w, r, "/", http.StatusFound)
}

func register (name string, display_name string, email string, pass string) int64 {
   password, salt := createHashFromPassword(pass)
   reg, err := db.Exec(`INSERT INTO users (name, display_name, salt, password, email) VALUES (?, ?, ?, ?, ?)`, name, display_name, salt, password, email)
   if err != nil {
      panic(err)
   }
   lastInsertID, _ := reg.LastInsertId()
   return lastInsertID
}

func createHashFromPassword(password string) (string, string) {
   bass := make([]byte, 14)
   _, err := io.ReadFull(rand.Reader, bass)
   if err != nil {
      fmt.Println("error:", err)
   }
   salt := base64.StdEncoding.EncodeToString(bass)
   converted, _ := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 16)
   return hex.EncodeToString(converted[:]), salt
}

func checkHashFromPassword(password string, salt string, hash string) bool {
   converted, _ := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 16)
   tmpHash := hex.EncodeToString(converted[:])
   if tmpHash == hash {
      return true
   }
   return false
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
   if !authenticated(w, r) {
      return
   }
   render(w, r, http.StatusOK, "timeline.html", nil)
}

func main(){
   host := "localhost"
   dbPortStr := "3306"
   dbPort , err := strconv.Atoi(dbPortStr)
   user := "root"
   dbName := "gotwitter"
   password :=  os.Getenv("DB_PASSWORD")
   db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local&parseTime=true", user, password, host, dbPort, dbName,))
   if err != nil { 
      //log.Fatalf("Failed to connect to DB: %s.", err.Error())
   }
   defer db.Close()

   secret := "hogehoge"
   store = sessions.NewCookieStore([]byte(secret)) 
   r := mux.NewRouter()
   l := r.Path("/login").Subrouter()
   l.Methods("GET").HandlerFunc(GetLogin)
   l.Methods("POST").HandlerFunc(PostLogin)
   r.Path("/logout").Methods("GET").HandlerFunc(GetLogout)
   
   reg := r.Path("/register").Subrouter()
   reg.Methods("GET").HandlerFunc(GetRegister)
   reg.Methods("POST").HandlerFunc(PostRegister)

   r.HandleFunc("/", GetIndex).Methods("GET")

   http.ListenAndServe(":8080", r)
}

