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

func processor(w http.ResponseWriter, r *http.Request){
	if r.Method != "Post"{
		http.Redirect(w,r,"/",http.StatusSeeOther)
		return
	}

	startpoint := r.FormValue("fstart")
	endpoint := r.FormValue("fend")
}

func main(){
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/",http.StripPrefix("/css", fs))
	http.HandleFunc("/",handleFunc)
	http.HandleFunc("/process",processor)
	http.ListenAndServe(":9999",nil)
}