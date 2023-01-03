package main



import (
//	"bufio"
	"flag"
	"fmt"
//	"io/ioutil"
        "github.com/rwcarlsen/goexif/exif"
//        "path"
	"log"
	"os"
	"os/exec"
	"path/filepath"
//	"strings"
	"time"
//        "github.com/nfnt/resize"
//	"image/jpeg"
        // "runtime"
)

// Post struct
type Post struct {
	Layout   string
	Title    string
	Date     string
	Category string
	Image    string
}

// Width of output image
var Width string = "512"

func main() {

	flag.Parse()
	values := flag.Args()[0]
	log.Println(values)

        now := fmt.Sprintf("%s", time.Now().Format(time.RFC3339)[0:10])
        log.Println(now)


        f, err := os.Open(values)
		if err != nil {
			panic(err)
		}
        defer f.Close()
		x, err := exif.Decode(f)
		if err != nil {
			panic(err)
		}
		tm, _ := x.DateTime()
		fmt.Println(tm.Date())


                /*
	fmt.Println("Enter Post title: ")

	reader := bufio.NewReader(os.Stdin)
	title, err := reader.ReadString('\n')
	title = title[0 : len(title)-1] //trim \n
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Enter Post date: ")

	date, err := reader.ReadString('\n')
	date = fmt.Sprintf("%s %s", date[0:len(date)-1], now) //trim \n
	if err != nil {
		log.Fatal(err)
	}
        */


        _, fn := filepath.Split(values)
        
        title := fn

        y,m,d := tm.Date()
        date := fmt.Sprintf("%04d-%02d-%02d",y,m,d)
        log.Println(date)
        
        _, file := filepath.Split(values)

	post := Post{
		Layout:   "post",
		Title:    title,
		Date:     date,
		Category: "kof archive",
		Image:    file,
	}

	log.Println(post)
        filename := filepath.Join("/", "home", "kof", "src", "k0f.github.io", "_posts", fmt.Sprintf("%s-%s.markdown", now, fn))
	log.Printf("Creating %s\n",filename)
        f, e := os.Create(filename)

	if e != nil {
		log.Fatal(e)
	}

	defer f.Close()

	//copy(values, filepath.Join("/home/kof/src/k0f.github.io/assets"), Width)

	_, err2 := f.WriteString(postToString(post))

	if err2 != nil {
		log.Fatal(err2)
	}


	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("convert %s -resize %s /home/kof/src/k0f.github.io/assets/%s; cd /home/kof/src/k0f.github.io; git add .; git commit -am \"%s\"; git push", file, Width, file, title)).Output()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Command Successfully Executed")
	output := string(out[:])
	fmt.Println(output)
}

func postToString(_post Post) string {
	return fmt.Sprintf("---\nlayout:\t%s\ntitle:\t\"%s\"\ndate:\t%s\ncategories:\t%s\n---\n\n![Image Alt](https://k0f.github.io/assets/%s)", _post.Layout, _post.Title, _post.Date, _post.Category, _post.Image)
}
