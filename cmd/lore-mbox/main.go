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
	m := flag.String("m", "html", "Mode: \"html\" or \"json\"")
	n := flag.String("n", "test", "Name")
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

	content := ""
	file := *n
	switch *m {
	case "html":
		lines := thread.ParseHTML(*c, *j)
		content = strings.Join(lines, "\n")
		file += ".html"

	case "json":
		data, err := thread.ParseJSON()
		if err != nil {
			log.Fatal(err)
		}
		content = string(data)
		file += ".json"

	default:
		err = fmt.Errorf("no such mode: %v", *m)
		log.Fatal(err)
	}

	err = os.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
