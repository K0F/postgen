package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Post struct {
	Layout   string
	Title    string
	Date     string
	Category string
	Image    string
}

var Width string = "512"

func main() {
	flag.Parse()
	
	// Prevent panic if no argument is passed
	if len(flag.Args()) < 1 {
		log.Fatal("Error: Please provide a path to an image file.\nUsage: postgen <path_to_image>")
	}
	
	values := flag.Args()[0]
	log.Println("Input path:", values)

	now := time.Now().Format("2006-01-02")
	log.Println("Current date:", now)

	_, file := filepath.Split(values)
	
	// Clean the title: remove extension and replace colons/spaces to protect Jekyll
	title := strings.Split(file, ".")[0]
	title = strings.ReplaceAll(title, ":", "-") 

	post := Post{
		Layout:   "post",
		Title:    title,
		Date:     now,
		Category: "kof archive",
		Image:    file,
	}

	// Sanitize output filename to prevent Jekyll URI scheme errors
	safeFn := strings.ReplaceAll(file, ":", "-")
	filename := filepath.Join("/", "home", "kof", "src", "k0f.github.io", "_posts", fmt.Sprintf("%s-%s.markdown", now, safeFn))
	log.Printf("Creating %s\n", filename)
	
	f, e := os.Create(filename)
	if e != nil {
		log.Fatal(e)
	}
	defer f.Close()

	_, err2 := f.WriteString(postToString(post))
	if err2 != nil {
		log.Fatal(err2)
	}

	// FIX: Use 'values' for source file so convert knows where to find it. 
	// FIX: Changed git commit -am to git commit -m
	shellCmd := fmt.Sprintf("convert \"%s\" -resize %s /home/kof/src/k0f.github.io/assets/\"%s\" && cd /home/kof/src/k0f.github.io && git add . && git commit -m \"archive: %s\" && git push", values, Width, file, title)
	
	log.Println("Executing pipeline...")
	out, err := exec.Command("/bin/sh", "-c", shellCmd).CombinedOutput() // CombinedOutput catches stderr too

	if err != nil {
		log.Printf("Shell pipeline failed: %s\n", err)
		log.Fatalf("Output: %s\n", string(out))
	}

	fmt.Println("Command Successfully Executed")
	fmt.Println(string(out))
}

// FIX: Replaced tabs (\t) with standard spaces to satisfy YAML specifications
func postToString(_post Post) string {
	return fmt.Sprintf("---\nlayout: %s\ntitle: \"%s\"\ndate: %s\ncategories: %s\n---\n\n![Image Alt](https://k0f.github.io/assets/%s)", 
		_post.Layout, _post.Title, _post.Date, _post.Category, _post.Image)
}