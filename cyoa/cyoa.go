package cyoa

import (
	"html/template"
	"net/http"
	"strings"
)

// Option ..
type Option struct {
	Text string
	Arc  string
}

// Arc ..
type Arc struct {
	Title   string
	Story   []string
	Options []Option
}

// Story ..
type Story map[string]Arc

// Story ..
type StoryHandler struct {
	Story Story
	Start string
}

// NewStoryHandler ..
func NewStoryHandler(s Story, start string) StoryHandler {
	return StoryHandler{
		Story: s,
		Start: start,
	}
}

// ServeHTTP ..
func (h StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
