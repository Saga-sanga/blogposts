package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Saga-sanga/blogposts"
	"github.com/Saga-sanga/blogposts/blogrenderer"
)

func main() {
	posts, err := blogposts.NewPostsFromFS(os.DirFS("posts"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(posts)

	postrenderer, err := blogrenderer.NewPostRenderer()
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	postrenderer.RenderIndex(&buf, posts)
	createHTMLFile(buf.String(), "index")

	for _, post := range posts {
		buf.Reset()
		postrenderer.Render(&buf, post)
		createHTMLFile(buf.String(), post.SanitisedTitle())
	}
}

func createHTMLFile(html, name string) {
	// get app directory path
	appDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// create directory if it does not exist
	htmlDir := filepath.Join(appDir, "/html")
	dir := filepath.Dir(htmlDir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(htmlDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// create new file
	filename := filepath.Join(htmlDir, fmt.Sprintf("/%s.html", name))
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// write data to file
	_, err = file.Write([]byte(html))
	if err != nil {
		log.Fatal(err)
	}
}
