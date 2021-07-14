package main

import (
	"log"
	"os"
	"strings"

	"github.com/FreeFlyingSheep/go-lore-mbox/pkg/lore"
	"github.com/FreeFlyingSheep/go-lore-mbox/pkg/mbox"
)

func main() {
	url := "https://lore.kernel.org/linux-arch/CAK8P3a2Qu_BUcGFpgktXOwsomuhN6aje6mB6EwTka0GBaoL4hw@mail.gmail.com/t.mbox.gz"

	url, err := lore.Parse(url)
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

	lines := thread.Parse("assets/style.css", "assets/tools.js")
	content := strings.Join(lines, "\n")
	os.WriteFile("test.html", []byte(content), os.ModePerm)
}
