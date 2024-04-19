package src

import(
	"html/template"
	"net/http"
)

var temp *template.Template

func init(){
	temp = template.Must(template.ParseGlob("html/*.html"))
}

func handleFunc(w http.ResponseWriter, r *http.Request){
	temp.ExecuteTemplate(w, "mainpage.html",nil)
}

func main(){
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/",http.StripPrefix("/css", fs))
	http.HandleFunc("/",handleFunc)
	http.ListenAndServe(":9999",nil)
}