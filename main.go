package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var dir_flag = flag.String("dir", "", "the directory for the files to be renamed")
var pat_flag = flag.String("pattern", "", "the regex for the target files ( single quoted )")

func processTargetDir(dir string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", nil
	}

	if dir == "" {
		return cwd, nil
	}

	if filepath.IsAbs(dir) {
		return dir, nil
	} else {
		return filepath.Join(cwd, dir), nil
	}

}

// ========= The logic for renaming files. Customize as needed. =============
func NewFileName(origFileName string) string {
	pattern := `(.+)_(\d\d\d)\.txt`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(origFileName)
	part1, part2 := matches[1], matches[2]
	return fmt.Sprintf("%s - %s.txt", part2, part1)
}

// ==========================================================================

func main() {
	flag.Parse()
	dir := *dir_flag

	targetPattern := *pat_flag
	targetDir, err := processTargetDir(dir)
	if err != nil {
		log.Fatalf("Processing targetDir failed: %s", err)
	}

	var walkFn fs.WalkDirFunc
	walkFn = func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		origFileName := d.Name()
		matched, e := regexp.MatchString(targetPattern, origFileName)
		if e != nil || !matched {
			return nil
		}
		newFileName := NewFileName(origFileName)
		newPath := filepath.Join(filepath.Dir(path), newFileName)
		if err := os.Rename(path, newPath); err != nil {
			log.Printf("Failed to rename %s\n", path)
			return nil // continue walking to other paths
		}
		fmt.Printf("Renamed %s to %s\n", path, newPath)
		return nil
	}

	err = filepath.WalkDir(targetDir, walkFn)
	if err != nil {
		log.Fatalf("error walking file tree: %s", err)
	}

}
