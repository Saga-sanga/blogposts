package blogrenderer

import (
	"bytes"
	"embed"
	"html/template"
	"io"

	"github.com/Saga-sanga/blogposts"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	//go:embed "templates/*"
	postTemplates embed.FS
)

type PostRenderer struct {
	templ    *template.Template
	markdown goldmark.Markdown
}

func NewPostRenderer() (*PostRenderer, error) {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)

	return &PostRenderer{templ: templ, markdown: markdown}, nil
}

func (r *PostRenderer) Render(w io.Writer, p blogposts.Post) error {
	data, err := newPostVM(p, r)
	if err != nil {
		return err
	}

	return r.templ.ExecuteTemplate(w, "blog.gohtml", data)
}

func (r *PostRenderer) RenderIndex(w io.Writer, posts []blogposts.Post) error {
	return r.templ.ExecuteTemplate(w, "index.gohtml", posts)
}

type postViewModel struct {
	blogposts.Post
	HTMLBody template.HTML
}

func newPostVM(p blogposts.Post, r *PostRenderer) (postViewModel, error) {
	var buf bytes.Buffer

	if err := r.markdown.Convert([]byte(p.Body), &buf); err != nil {
		return postViewModel{}, err
	}

	vm := postViewModel{Post: p}
	vm.HTMLBody = template.HTML(buf.String())
	return vm, nil
}
