package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
)

var ignore = regexp.MustCompile(`\.git|build|testdata|3rdparty|packages`)
var verbose bool

func adddirs(pkg, path string, dirs []string) []string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if fi, err := f.Readdir(0); err != nil {
		panic(err)
	} else {
		for _, f := range fi {
			if !f.IsDir() || ignore.MatchString(f.Name()) {
				continue
			}
			pkg2 := fmt.Sprintf("%s/%s", pkg, f.Name())
			dirs = append(dirs, pkg2)
			dirs = adddirs(pkg2, fmt.Sprintf("%s/%s", path, f.Name()), dirs)
		}
	}
	return dirs
}

func copy(a, b string) {
	f1, err := os.Open(a)
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	f2, err := os.Create(b)
	if err != nil {
		panic(err)
	}
	defer f2.Close()
	if verbose {
		fmt.Printf("Copying %s to %s\n", f1.Name(), f2.Name())
	}
	if _, err := io.Copy(f2, f1); err != nil {
		panic(err)
	}
}

func buildPeg(pegFile, ignore, testfile string, extra ...string) {
	args := []string{"-peg=" + pegFile, "-notest", "-ignore", ignore, "-testfile", testfile}
	args = append(args, extra...)
	c := exec.Command("./parser_exe", args...)
	if verbose {
		fmt.Println(c.Args)
	}
	if b, err := c.CombinedOutput(); err != nil {
		panic(fmt.Sprint(string(b), err))
	} else if len(b) != 0 {
		fmt.Println(string(b))
	}
	outname := filepath.Base(pegFile)
	outname = outname[:len(outname)-4]
	copy(fmt.Sprintf("%s/%s.go", outname, outname), pegFile[:len(pegFile)-4]+".go")
	copy(fmt.Sprintf("%s/%s_test.go", outname, outname), pegFile[:len(pegFile)-4]+"_test.go")
}

func readthread(r io.Reader, out chan string) {
	buf := make([]byte, 2048)
	for {
		if n, err := r.Read(buf); err != nil && err != io.EOF {
			panic(err)
		} else if err == io.EOF {
			break
		} else {
			out <- string(buf[:n])
		}
	}
	close(out)
}

func main() {
	flag.BoolVar(&verbose, "v", verbose, "Verbose output")
	flag.Parse()
	var extension = ""
	if runtime.GOOS == "windows" {
		extension += ".exe"
	}
	c := exec.Command("go", "build", "-o", "parser_exe"+extension, "github.com/quarnster/parser/exe")
	if verbose {
		fmt.Println(c.Args)
	}
	if b, err := c.CombinedOutput(); err != nil {
		panic(fmt.Sprint(string(b), err))
	} else if len(b) != 0 {
		fmt.Println(string(b))
	}
	buildPeg("../backend/loaders/json/json.peg", "JsonFile,Values,Value,Null,Dictionary,Array,KeyValuePairs,KeyValuePair,QuotedText,Text,Integer,Float,Boolean,Spacing,Comment", "testdata/Default (OSX).sublime-keymap")
	buildPeg("../backend/loaders/plist/plist.peg", "Spacing,KeyValuePair,KeyTag,StringTag,Value,Values,PlistFile,Plist", "../../../3rdparty/bundles/c.tmbundle/Syntaxes/C.plist", "-dumptree")

	c = exec.Command("go", "run", "python.go")
	if verbose {
		fmt.Println(c.Args)
	}
	if b, err := c.CombinedOutput(); err != nil {
		panic(fmt.Sprint(string(b), err))
	} else if len(b) != 0 {
		fmt.Println(string(b))
	}

	tests := []string{"test"}
	if verbose {
		tests = append(tests, "-v")
	}
	//	tests = append(tests, "lime")
	tests = adddirs("lime", "..", tests)
	c = exec.Command("go", tests...)
	r, err := c.StdoutPipe()
	if err != nil {
		panic(err)
	}
	r2, err := c.StderrPipe()
	if err != nil {
		panic(err)
	}
	if err := c.Start(); err != nil {
		panic(err)
	}
	sc := make(chan string)
	ec := make(chan string)
	go readthread(r, sc)
	go readthread(r2, ec)
	so, eo := true, true
	for so && eo {
		line := ""
		select {
		case line, so = <-sc:
		case line, eo = <-ec:
		}
		fmt.Print(line)
	}
}
