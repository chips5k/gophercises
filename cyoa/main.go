package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type option struct {
	Text string
	Arc  string
}

type arc struct {
	Title   string
	Story   []string
	Options []option
}

type story map[string]arc

type handler struct {
	Story story
	Start string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arc := strings.TrimLeft(r.URL.Path, "/")

	if arc == "" {
		arc = "intro"
	}

	content, ok := h.Story[arc]

	if !ok {
		http.Error(w, "Not found", 404)
		return
	}

	t, err := template.New("main").Parse(`
		<!doctype html>
		<html>
			<body>
				<h1>{{.Title}}</h1>	
				{{ range .Story }}
					<p>{{ . }}
				{{ end }}
				<hr />
				<ul>
				{{ range .Options }}
					<li><a href="/{{ .Arc }}">{{ .Text }}</a></li>
				{{ end }}
				</ul>
			</body>
		</html>
	`)

	if err != nil {
		http.Error(w, "Something broke!", 500)
		return
	}

	t.Execute(w, content)
}

func main() {

	path := flag.String("json", "gopher.json", "path of a json story file")
	start := flag.String("start", "intro", "key of the first story arc")
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %s, error: %v", *path, err)
		os.Exit(1)
	}

	bb, err := ioutil.ReadAll(f)
	if err != nil {

		fmt.Fprintf(os.Stderr, "Failed to read file: %v", err)
		os.Exit(1)
	}

	s := make(story)
	err = json.Unmarshal(bb, &s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal json: %v", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	h := handler{
		Story: s,
		Start: *start,
	}

	mux.Handle("/", h)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
