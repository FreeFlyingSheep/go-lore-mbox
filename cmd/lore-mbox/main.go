package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/FreeFlyingSheep/go-lore-mbox/pkg/lore"
	"github.com/FreeFlyingSheep/go-lore-mbox/pkg/mbox"
)

func main() {
	m := flag.String("m", "html", "Mode: \"html\" or \"json\" or \"patch\"")
	o := flag.String("o", "test", "Output filename or directory")
	u := flag.String("u", "", "https://lore.kernel.org/xxx/xxx")
	c := flag.String("c", "assets/style.css", "CSS file, only works in html mode")
	j := flag.String("j", "assets/tools.js", "JS file, only works in html mode")
	flag.Parse()

	url, err := lore.Parse(*u)
	if err != nil {
		log.Fatal(err)
	}

	data, err := lore.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	messages, err := mbox.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	thread, err := mbox.Create("test", messages)
	if err != nil {
		log.Fatal(err)
	}

	contents := []string{}
	files := []string{}
	switch *m {
	case "html":
		lines := thread.ParseHTML(*c, *j)
		contents = append(contents, strings.Join(lines, "\n"))
		files = append(files, *o+".html")

	case "json":
		data, err := thread.ParseJSON()
		if err != nil {
			log.Fatal(err)
		}
		contents = append(contents, string(data))
		files = append(files, *o+".json")

	case "patch":
		err := os.MkdirAll(*o, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		data, err := thread.ParsePatch()
		if err != nil {
			log.Fatal(err)
		}
		for _, d := range data {
			contents = append(contents, d[1])
			files = append(files, *o+"/"+d[0]+".patch")
		}

	default:
		err = fmt.Errorf("no such mode: %v", *m)
		log.Fatal(err)
	}

	for i, file := range files {
		err = os.WriteFile(file, []byte(contents[i]), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
