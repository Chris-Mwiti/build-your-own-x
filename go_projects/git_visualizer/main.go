package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

//scanGitFolders returns a list of subfolders of 'folder' ending with '.git'
// Returns the base folder of the repo, the .git folder parent
// Recursively searches in the subfolder by passing an existing 'folders' slice.
func scanGitFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)
	f.Close()

	if err != nil {
		log.Fatal(err)
	}

	var path string

	//Recursively go through the folder and check whether it contains any .git files and if replace the current path with the found git path
	for _, file := range files {
		if file.IsDir() {
			path = folder + "/" + file.Name()
			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}

			//Omit both the vendor files and node_modules files due to size
			if file.Name() == "vendor" || file.Name() == "node_modules" {
				continue
			}

			//Recursive function
			folders = scanGitFolders(folders, path)
		}
	}
	
	return folders
}

func recursiveScanFolder(folder string) []string {
	return scanGitFolders(make([]string, 0), folder)
}

//returns the file path of where exactly the slice content should be stored
func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + "/.gogitlocalstatus"

	return dotFile

}


//parse the existing repos stored in the file to a slice
//add the new items to the slice, without adding duplicates
//store the slice to the file, overwriting the existing content

//1. Open the file
func openFile(filepath string) *os.File {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0755)

	if err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			_, err = os.Create(filepath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	return f
}

//2. Parse existing lines in the file into a slice
func parseFileLinesToSlice(filepath string) []string {
	f := openFile(filepath)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}
	
	return lines
}

//3. Join existing slices with new slice removing duplicates
func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !sliceContains(existing, i) {
			existing = append(existing, i)
		}
	}
	
	return existing
}

//utility func to check whether already existing content exists
func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	
	return false
}

//dumpstringsslice to file writes content to the file in path 'filepath' (overwriting existing content)
func dumpStringsSliceToFile(repos []string, filepath string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(filepath, []byte(content), 0755)
}

func addNewSliceElementsToFile(filePath string, newRepos []string) {
	existingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, existingRepos)
	dumpStringsSliceToFile(repos, filePath)
}

// crawls the given path and its subfolders
// searching for git repositories
func scan(folder string){
	fmt.Printf("Found folders:\n\n")
	repositories := recursiveScanFolder(folder)
	filepath := getDotFilePath()
	addNewSliceElementsToFile(filepath, repositories)
	fmt.Printf("\n\nSuccessfully added \n\n")
	print("scan")
}



func main(){
	var folder string
	var email string

	flag.StringVar(&folder, "add", "", "add a new folder to scan for git repositories")
	flag.StringVar(&email, "email", "your@email.com", "the email to scan")
	flag.Parse()

	if folder != "" {
		scan(folder)
		return
	}

	stats(email)
}