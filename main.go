package main

import (
	"log"
	"os"
)

func main() {
	defer func() {
		if v := recover(); v != nil {
			log.Fatal(v)
		}
	}()
	log.SetFlags(0)
	if len(os.Args) < 3 || (os.Args[1] != "verify" && os.Args[1] != "format") {
		log.Fatalf("Usage: %s ( verify | format ) <filename>...", os.Args[0])
	}
	for _, filename := range os.Args[2:] {
		grammar := parse(filename)
		switch os.Args[1] {
		case "format":
			print(filename, grammar, "SourceFile")
		case "verify":
			Verify(grammar, "SourceFile")
		}
	}
}

func parse(filename string) map[string]*Production {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	return Parse(filename, file)
}

func print(filename string, grammar map[string]*Production, start string) {
	file, _ := os.Create(filename)
	defer file.Close()
	Print(grammar, start, file)
}
