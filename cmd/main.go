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

	for _, post := range posts {
		var buf bytes.Buffer
		postrenderer.Render(&buf, post)

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
		filename := filepath.Join(htmlDir, fmt.Sprintf("/%s.html", post.SanitisedTitle()))
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// write data to file
		_, err = file.Write([]byte(buf.String()))
		if err != nil {
			log.Fatal(err)
		}
	}
}
