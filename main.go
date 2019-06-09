package main

import (
	"bufio"
	"errors"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/thoas/go-funk"

	"github.com/bmatcuk/doublestar"
)

type runtimeInfo struct {
	inDir    string
	outDir   string
	template *template.Template
	files    []string
}

var runInfo runtimeInfo

func parseFlags() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	inDir := flag.String("in", dir, "input directory")
	outDir := flag.String("out", dir, "output directory")

	flag.Parse()

	if *inDir == dir && *outDir == dir {
		*outDir = filepath.Join(*outDir, "generated")
	}

	runInfo.inDir = *inDir
	runInfo.outDir = *outDir
}

func initializeTemplate() {
	funcMap := template.FuncMap{
		"minus": func(a, b int) int {
			return a - b
		},
		"plus": func(a, b int) int {
			return a + b
		},
		"minus64": func(a, b int64) int64 {
			return a - b
		},
		"plus64": func(a, b int64) int64 {
			return a + b
		},
		"uuid": func() string {
			id, err := uuid.NewV4()
			if err == nil {
				return id.String()
			}
			return ""
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}

	extensions := []string{".html", ".template"}
	runInfo.files = []string{}
	f, err := doublestar.Glob(filepath.Join(runInfo.inDir, "**"))
	if err != nil {
		log.Fatal(err)
	}
	runInfo.files = append(runInfo.files, funk.FilterString(f, func(s string) bool {
		fileInfo, err := os.Stat(s)
		if err != nil {
			log.Fatal(err)
		}
		isInOut, err := doublestar.PathMatch(filepath.Join(runInfo.outDir, "**"), s)
		if err != nil {
			log.Fatal(err)
		}

		return !isInOut && !fileInfo.IsDir() && funk.ContainsString(extensions, filepath.Ext(s))
	})...)

	runInfo.template = template.New("main").Funcs(funcMap)
	for i, fname := range runInfo.files {
		fname, err := filepath.Rel(runInfo.inDir, fname)
		if err != nil {
			log.Fatal(err)
		}
		fname = filepath.ToSlash(fname)

		b, err := ioutil.ReadFile(runInfo.files[i])
		if err != nil {
			log.Fatal(err)
		}
		s := string(b)

		runInfo.template, err = runInfo.template.New(fname).Parse(s)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func chewUp() {
	// .html files produce pages; do not parse output files
	mainFiles := funk.FilterString(runInfo.files, func(s string) bool {
		isInOut, err := doublestar.PathMatch(filepath.Join(runInfo.outDir, "**"), s)
		if err != nil {
			log.Fatal(err)
		}
		return strings.HasSuffix(s, ".html") && !isInOut
	})

	for _, fname := range mainFiles {
		fname, err := filepath.Rel(runInfo.inDir, fname)
		if err != nil {
			log.Fatal(err)
		}
		outPath := filepath.Join(runInfo.outDir, fname)
		log.Println("Processing", fname)

		os.MkdirAll(filepath.Dir(outPath), os.ModePerm)

		f, err := os.Create(outPath)
		if err != nil {
			log.Fatal(err)
		}

		w := bufio.NewWriter(f)

		err = runInfo.template.ExecuteTemplate(w, fname, struct{}{})
		if err != nil {
			f.Close()
			log.Fatal(err)
		}
		w.Flush()
		f.Close()
	}
}

func main() {
	parseFlags()
	log.Println("Input directory is", runInfo.inDir)
	log.Println("Output directory is", runInfo.outDir)
	initializeTemplate()
	chewUp()
}
