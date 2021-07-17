package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/FreeFlyingSheep/go-lore-mbox/pkg/lore"
	"github.com/FreeFlyingSheep/go-lore-mbox/pkg/mbox"
)

func main() {
	u := flag.String("u", "", "https://lore.kernel.org/xxx/xxx")
	c := flag.String("c", "assets/style.css", "css file")
	j := flag.String("j", "assets/tools.js", "js file")
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

	lines := thread.Parse(*c, *j)
	content := strings.Join(lines, "\n")
	err = os.WriteFile("test.html", []byte(content), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
