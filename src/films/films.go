package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sort"
	"text/template"
	"os"
	"apicaller"
)
//Struct holding a Movie (Film) information
type Film struct {
	Title string
	Note  string
}

func (p Film) String() string {
	return fmt.Sprintf("%s: %s", p.Title, p.Note)
}

// ByNote implements sort.Interface for []Film based on
// the Note field.
type ByNote []Film

func (a ByNote) Len() int           { return len(a) }
func (a ByNote) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNote) Less(i, j int) bool { return a[i].Note < a[j].Note }

func findnote(c chan Film, element []string, i int, apimax chan bool){
	fmt.Println("Doing", element[1])
	str := strings.Replace(element[1]," ","+",-1)
	m, err := apicaller.Callapisem("http://www.omdbapi.com/?t="+str+"&y=&plot=short&r=json", apimax)
	if (err != nil){
		f := Film{element[1],"Error"}
		c <- f
		return
	}
	if m["Response"] == "True" {
		f := Film{element[1],m["imdbRating"].(string)}
		c <- f
	} else {
		f := Film{element[1],m["Error"].(string)}
		c <- f
	}
}

func reply(w http.ResponseWriter, s string){
	t, _ := template.ParseFiles(os.Getenv("GOPATH")+"/lib/films/result.html")
	t.ExecuteTemplate(w, "Body", s)
}

func handler2(w http.ResponseWriter, r *http.Request) {
	src := r.FormValue("bodyform")
	re := regexp.MustCompile(`<h4><[^>]*>([^<]*)`) 
	res := re.FindAllStringSubmatch(src, -1)
	if len(res) <= 6{
		reply(w, fmt.Sprintf("Invalid input, %s", res))
		return
	}
	finalmap := make([]Film, len(res)-6)
	c := make(chan Film)
	//This defines the max number of parralel requests the API can withstand
	apimax := make(chan bool, 20)
	for i,element := range res[:len(res)-6] {
		go findnote(c,element,i,apimax)
	}
	for i:=0; i< len(res)-6; i++ {
		f := <-c
		finalmap[i].Title = f.Title
		finalmap[i].Note = f.Note
	}
	sort.Sort(ByNote(finalmap))
	finalstring := ""
	for j := len(finalmap)-1; j >= 0; j-- {
		finalstring = finalstring + fmt.Sprintf("<p>%s: %s</p>", finalmap[j].Title, finalmap[j].Note)
	}

	reply(w, finalstring)
	
}
func main() {
	
	http.Handle("/", http.FileServer(http.Dir(os.Getenv("GOPATH")+"/lib/films/")))
	http.HandleFunc("/process/", handler2)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(os.Getenv("GOPATH")+"/lib/films/css"))))
	http.ListenAndServe(":8080", nil)
}
