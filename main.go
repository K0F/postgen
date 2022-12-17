package main



import (
	"bufio"
	"flag"
	"fmt"
	//"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
        "github.com/nfnt/resize"
	"image/jpeg"
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
var Width uint = 512

func copy(src, dst string, width uint) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()


        // decode jpeg into image.Image
	img, err := jpeg.Decode(source)
	if err != nil {
		log.Fatal(err)
	}
        m := resize.Resize(width, 0, img, resize.Lanczos3)


	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	
        //nBytes, err := io.Copy(destination, source)
        
        jpeg.Encode(destination, m, nil)

        return err
}

func main() {

	flag.Parse()
	values := flag.Args()[0]
	log.Println(values)

	now := strings.Split(fmt.Sprintf("%s", time.Now().Format(time.UnixDate)), " ")[3]

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

	post := Post{
		Layout:   "post",
		Title:    title,
		Date:     date,
		Category: "kof archive",
		Image:    values,
	}

	filename := filepath.Join("/", "home", "kof", "src", "k0f.github.io", "_posts", fmt.Sprintf("%s-%s.markdown", strings.Split(post.Date, " ")[0], strings.Replace(post.Title, " ", "-", -1)))
	f, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	copy(values, filepath.Join("/home/kof/src/k0f.github.io/assets"), Width)

	_, err2 := f.WriteString(postToString(post))

	if err2 != nil {
		log.Fatal(err2)
	}

	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd /home/kof/src/k0f.github.io; git add .; git commit -am \"%s\";git push", title)).Output()

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
