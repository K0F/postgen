package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	postsDir := "/home/kof/src/k0f.github.io/_posts"
	assetsDir := "/home/kof/src/k0f.github.io/assets"

	files, err := os.ReadDir(postsDir)
	if err != nil {
		log.Fatal(err)
	}

	reYYYY := regexp.MustCompile(`(20\d{2})[_-]?(0[1-9]|1[0-2])[_-]?(0[1-9]|[12]\d|3[01])`)
	reYY := regexp.MustCompile(`(?:IMG[_-]?)(\d{2})(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])`)
	imageRegex := regexp.MustCompile(`assets/([^\s\)\"]+)`)

	processed := 0
	renamed := 0
	skipped := 0

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".markdown") {
			continue
		}

		oldPath := filepath.Join(postsDir, file.Name())
		contentBytes, err := os.ReadFile(oldPath)
		if err != nil {
			log.Printf("[ERROR] Cannot read %s: %v\n", file.Name(), err)
			continue
		}

		content := string(contentBytes)
		imgMatches := imageRegex.FindStringSubmatch(content)
		if len(imgMatches) < 2 {
			log.Printf("[SKIP] %s -> No asset image match\n", file.Name())
			skipped++
			continue
		}

		imgName := imgMatches[1]
		imgPath := filepath.Join(assetsDir, imgName)

		realDate := ""
		matchesYYYY := reYYYY.FindStringSubmatch(imgName)
		if len(matchesYYYY) == 4 {
			realDate = fmt.Sprintf("%s-%s-%s", matchesYYYY[1], matchesYYYY[2], matchesYYYY[3])
		} else {
			matchesYY := reYY.FindStringSubmatch(imgName)
			if len(matchesYY) == 4 {
				realDate = fmt.Sprintf("20%s-%s-%s", matchesYY[1], matchesYY[2], matchesYY[3])
			} else {
				out, err := exec.Command("exiftool", "-DateTimeOriginal", "-d", "%Y-%m-%d", "-s3", imgPath).Output()
				if err == nil {
					d := strings.TrimSpace(string(out))
					if len(d) == 10 {
						realDate = d
					}
				}
			}
		}

		if realDate == "" {
			if fi, err := os.Stat(imgPath); err == nil {
				realDate = fi.ModTime().Format("2006-01-02")
			}
		}

		if realDate == "" {
			log.Printf("[SKIP] %s -> Could not determine date\n", file.Name())
			skipped++
			continue
		}

		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "date:") {
				lines[i] = fmt.Sprintf("date: %s", realDate)
				break
			}
		}
		newContent := strings.Join(lines, "\n")

		parts := strings.SplitN(file.Name(), "-", 4)
		var baseName string
		if len(parts) == 4 && len(parts[0]) == 4 && len(parts[1]) == 2 && len(parts[2]) == 2 {
			baseName = parts[3]
		} else {
			baseName = file.Name()
		}

		newFileName := fmt.Sprintf("%s-%s", realDate, baseName)
		newPath := filepath.Join(postsDir, newFileName)

		if err := os.WriteFile(oldPath, []byte(newContent), 0644); err != nil {
			log.Printf("[ERROR] Writing %s: %v\n", oldPath, err)
			continue
		}

		if oldPath != newPath {
			if err := os.Rename(oldPath, newPath); err != nil {
				log.Printf("[ERROR] Renaming %s to %s: %v\n", oldPath, newPath, err)
			} else {
				fmt.Printf("[RENAMED] %s -> %s\n", file.Name(), newFileName)
				renamed++
			}
		} else {
			fmt.Printf("[UPDATED] %s\n", file.Name())
		}
		processed++
	}

	fmt.Printf("\nDone. Processed: %d, Renamed: %d, Skipped: %d\n", processed, renamed, skipped)
}