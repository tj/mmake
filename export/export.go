package export

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Export(makefilePath string, b *bytes.Reader) (string, error) {
	path, err := filepath.Abs(makefilePath)
	if err != nil {
		log.Fatalf("Cannot find full path of %s", makefilePath)
	}
	fmt.Printf("Path: %s", path)

	includes := IncludedMakefilePaths(path)
	return fmt.Sprintf("%+v", includes), nil
}

func ProcessMakefile(path string, b *bytes.Reader, includes []string) []string {
	// TODO: figure out how to save this info and append to output string
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasPrefix(txt, "include ") {
			// recursive call process_makefile
			line := strings.SplitN(txt, " ", 2)
			// lookup line in include paths
			ProcessMakefile(line[1], bytes.NewReader(b), includes)
		} else {
		}
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}

func IncludedMakefilePaths(path string) []string {
	stubMakefile := fmt.Sprintf(`
include %s

mmake_includes:
	@echo $(MAKEFILE_LIST)
`, path)

	content := []byte(stubMakefile)
	tmpfile, err := ioutil.TempFile("", "mmake")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	// Use tempfile
	out, err := exec.Command("make", "-f", tmpfile.Name(), "mmake_includes").Output()
	if err != nil {
		log.Fatal(err)
	}
	files := strings.Split(string(out), " ")

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return files[1:]
}
