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
	var buf bytes.Buffer

	if err := r.ConvertToMarkdown(&buf, p.Body); err != nil {
		return err
	}

	data := struct {
		blogposts.Post
		HTMLBody template.HTML
	}{
		Post:     p,
		HTMLBody: template.HTML(buf.String()),
	}

	if err := r.templ.ExecuteTemplate(w, "blog.gohtml", data); err != nil {
		return err
	}

	return nil
}

func (r *PostRenderer) ConvertToMarkdown(w io.Writer, mdContent string) error {
	if err := r.markdown.Convert([]byte(mdContent), w); err != nil {
		return err
	}

	return nil
}
