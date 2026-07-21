package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

func extractDate(filePath string, customDate string) string {
	if customDate != "" {
		return customDate
	}

	_, file := filepath.Split(filePath)

	reYYYY := regexp.MustCompile(`(20\d{2})[_-]?(0[1-9]|1[0-2])[_-]?(0[1-9]|[12]\d|3[01])`)
	matchesYYYY := reYYYY.FindStringSubmatch(file)
	if len(matchesYYYY) == 4 {
		return fmt.Sprintf("%s-%s-%s", matchesYYYY[1], matchesYYYY[2], matchesYYYY[3])
	}

	reYY := regexp.MustCompile(`(?:IMG[_-]?)(\d{2})(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])`)
	matchesYY := reYY.FindStringSubmatch(file)
	if len(matchesYY) == 4 {
		return fmt.Sprintf("20%s-%s-%s", matchesYY[1], matchesYY[2], matchesYY[3])
	}

	out, err := exec.Command("exiftool", "-DateTimeOriginal", "-d", "%Y-%m-%d", "-s3", filePath).Output()
	if err == nil {
		exifDate := strings.TrimSpace(string(out))
		if len(exifDate) == 10 {
			return exifDate
		}
	}

	if fi, err := os.Stat(filePath); err == nil {
		return fi.ModTime().Format("2006-01-02")
	}

	return time.Now().Format("2006-01-02")
}

func main() {
	var customDate string
	flag.StringVar(&customDate, "date", "", "Custom date (YYYY-MM-DD)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Error: Please provide a path to an image file.\nUsage: postgen [-date YYYY-MM-DD] <path_to_image>")
	}

	values := flag.Args()[0]
	log.Println("Input path:", values)

	postDate := extractDate(values, customDate)
	log.Println("Post date:", postDate)

	_, file := filepath.Split(values)

	title := strings.Split(file, ".")[0]
	title = strings.ReplaceAll(title, ":", "-")

	post := Post{
		Layout:   "post",
		Title:    title,
		Date:     postDate,
		Category: "kof archive",
		Image:    file,
	}

	safeFn := strings.ReplaceAll(file, ":", "-")
	filename := filepath.Join("/", "home", "kof", "src", "k0f.github.io", "_posts", fmt.Sprintf("%s-%s.markdown", postDate, safeFn))
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

	shellCmd := fmt.Sprintf("convert \"%s\" -resize %s /home/kof/src/k0f.github.io/assets/\"%s\" && cd /home/kof/src/k0f.github.io && git add . && git commit -m \"archive: %s\" && git push", values, Width, file, title)

	log.Println("Executing pipeline...")
	out, err := exec.Command("/bin/sh", "-c", shellCmd).CombinedOutput()

	if err != nil {
		log.Printf("Shell pipeline failed: %s\n", err)
		log.Fatalf("Output: %s\n", string(out))
	}

	fmt.Println("Command Successfully Executed")
	fmt.Println(string(out))
}

func postToString(_post Post) string {
	return fmt.Sprintf("---\nlayout: %s\ntitle: \"%s\"\ndate: %s\ncategories: %s\n---\n\n![Image Alt](https://k0f.github.io/assets/%s)",
		_post.Layout, _post.Title, _post.Date, _post.Category, _post.Image)
}