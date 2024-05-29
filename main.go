package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
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

func generateFilename() string {
	var bytes = make([]byte, 4)
	rand.Read(bytes[:])
	return fmt.Sprintf("%#x", bytes)
}

var chromeFlags = []string{
	"--headless",
}

func generatePDFFromSource(source []byte) ([]byte, error) {
	tmpDir := os.TempDir()
	id := generateFilename()
	tmpFile := path.Join(tmpDir, id+".html")
	outFile := path.Join(tmpDir, id+".pdf")
	args := append(chromeFlags, fmt.Sprintf(`--print-pdf="%s"`, outFile), tmpFile)

	os.WriteFile(tmpFile, source, 0o600)

	out := bytes.NewBuffer([]byte{})

	cmd := exec.Command("chromium", args...)
	cmd.Stderr = out
	cmd.Stdout = out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(outFile)
}

func main() {
	source := bytes.NewBuffer([]byte{})
	pdfTemplate.Execute(source, invoiceAssigns{
		Styles: styles,
	})

	fmt.Println(source.String())
}
