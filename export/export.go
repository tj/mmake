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

	includes := IncludedMakefilePaths(path)
	exported := ProcessMakefile(path, includes, []string{})
	return fmt.Sprintf("%+v", strings.Join(exported, "\n")+"\n"), nil
}

func read(file string) ([]byte, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

func getFullIncludePath(partial string, includes []string) string {
	for _, fullPath := range includes {
		m := strings.Index(fullPath, partial)
		//fmt.Printf("Is match? %s %s %d\n", partial, fullPath, m)
		if m != -1 {
			return fullPath
		}
	}
	return ""
}

func ProcessMakefile(path string, includes []string, outputLines []string) []string {
	// TODO: figure out how to save this info and append to output string
	f, err := read(path)
	if err != nil {
		log.Fatalf("Unable to read %s with error: %s", path, err)
	}
	b := bytes.NewBuffer(f)
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasPrefix(txt, "include ") {
			outputLines = append(outputLines, fmt.Sprintf("#- start=%s", txt))
			// recursive call process_makefile
			line := strings.SplitN(txt, " ", 2)
			// lookup line in include paths
			fullPath := getFullIncludePath(line[1], includes)
			for _, l := range ProcessMakefile(fullPath, includes, []string{}) {
				outputLines = append(outputLines, l)
			}
			outputLines = append(outputLines, fmt.Sprintf("#- end=%s", txt))
		} else {
			outputLines = append(outputLines, txt)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return outputLines
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
	var result []string

	for i, f := range files {
		ft := strings.TrimSpace(f)
		if i != 0 && ft != "" {
			result = append(result, ft)
		}
	}

	return result
}
