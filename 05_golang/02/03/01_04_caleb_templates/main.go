package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	log.SetFlags(0)

	tpl, err := template.ParseFiles("tpl.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tpl.Execute(os.Stdout, 5*5)
	if err != nil {
		log.Fatalln(err)
	}
}


