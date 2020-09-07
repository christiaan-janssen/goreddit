package web

import (
	"html/template"
	"net/http"

	"github.com/christiaan-janssen/goreddit"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	h.Use(middleware.Logger)
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.ThreadList())
	})
	return h
}

type Handler struct {
	*chi.Mux

	store goreddit.Store
}

const threadListHTML = `
<h1>Threads</h1>
{{range .Threads}}
	<dt><strong>{{.Title}}</strong></dt>
	<dd>{{.Description}}</dd>
{{end}}
`

func (h *Handler) ThreadList() http.HandlerFunc {
	type data struct {
		Threads []goreddit.Thread
	}

	tmpl := template.Must(template.New("").Parse(threadListHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{Threads: tt})
	}
}

const threadCreateHTML = `
<h1>New Thread</h1>
<form action="/threads" method="POST">
  <table>
    <tr>
      <td>Title</td>
      <td><input type="text" name="title" /></td>
    </tr>
    <tr>
      <td>Description</td>
      <td><input type="text" name="description" /></td>
    </tr>
    </table>
    <button type="submit">Create thread</button>
</form>
`

func (h *Handler) ThreadCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(threadCreateHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}
