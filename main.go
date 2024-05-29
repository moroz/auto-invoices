package main

import (
	"fmt"
	"html/template"
	"os"
)

var styles template.CSS

func init() {
	contents, err := os.ReadFile("./templates/style.css")
	if err != nil {
		panic(err)
	}
	styles = template.CSS(contents)
}

var pdfTemplate = template.Must(template.ParseFiles("./templates/layout.html.tmpl"))

type invoiceAssigns struct {
	Styles template.CSS
}

func main() {
	fmt.Println(string(styles))
	pdfTemplate.Execute(os.Stdout, invoiceAssigns{
		Styles: styles,
	})
}
