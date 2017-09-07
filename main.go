package main

import(
   "github.com/gorilla/mux"
   "net/http"
   "html/template"
   "path"
)

func getTemplatePath(file string) string {
   return path.Join("contents", file)
}

func render(w http.ResponseWriter, r *http.Request, status int, file string, data interface{}){
   t := template.Must(template.New(file).ParseFiles(getTemplatePath(file)))
   w.WriteHeader(status)
   t.Execute(w, data)
}

func GetIndex(w http.ResponseWriter, r *http.Request){
   render(w, r, http.StatusOK, "timeline.html", nil)
}

func main(){
   // host := "localhost"
   // dbPortStr := "3306"
   // dbPort , err := strconv.Atoi(dbPortstr)
   // user := "root"
   // dbName := "gotwitter"

   r := mux.NewRouter()
   r.HandleFunc("/", GetIndex).Methods("GET")

   http.ListenAndServe(":8080", r)
}

