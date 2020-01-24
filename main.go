package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"text/template"
)

////Configurando template
var tpl *template.Template

var banco sync.Map

func init() {
	tpl = template.Must(template.ParseGlob("./Templates/*.html"))
}

///Config template fim

type dados struct {
	Data   []string `xml:"channel>item>pubDate"`
	Desc   []string `xml:"channel>item>description"`
	Titulo []string `xml:"channel>item>title"`
	Img    []string `xml:"channel>item>media:content"`
	Link   []string `xml:"channel>item>link"`
}

type dadosf struct {
	Data   string
	Desc   string
	Titulo string
	Img    int
	Link   string
}

type geral struct {
	Titulo string
	Dados  []dadosf
}

func hunter(w http.ResponseWriter, r *http.Request) {
	var source dados
	var sourcef []dadosf

	var urll string

	//fmt.Println("URL ? -> ", r.URL)
	urll = r.URL.String()

	resp, _ := http.Get("http://g1.globo.com/dynamo" + urll + "/rss2.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &source)

	for c := 0; c < 30; c++ {
		sourcef = append(sourcef, dadosf{Data: source.Data[c], Desc: source.Desc[c], Titulo: source.Titulo[c], Link: source.Link[c], Img: c})
		//fmt.Fprintf(w, "<h1>"+source.Titulo[c]+"</h1>"+source.Desc[c]+"\n\n")
	}
	//termo := r.URL.String()
	tpl.ExecuteTemplate(w, "index.html", geral{Titulo: urll, Dados: sourcef})

}

func erro(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Página Inexistente!   :(")
	tpl.ExecuteTemplate(w, "error.html", nil)
}

func main() {

	http.HandleFunc("/planeta-bizarro", hunter)
	http.HandleFunc("/carros", hunter)
	http.HandleFunc("/tecnologia", hunter)
	http.HandleFunc("/turismo-e-viagem", hunter)
	http.HandleFunc("/mundo", hunter)
	http.HandleFunc("/", erro)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
