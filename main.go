package main

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

var styles template.CSS

func init() {
	contents, err := os.ReadFile("./templates/style.css")
	if err != nil {
		panic(err)
	}
	styles = template.CSS(contents)
}

var pdfTemplate = template.Must(template.New("layout.html.tmpl").Funcs(template.FuncMap{
	"format_date": func(date time.Time) string {
		return date.Format("02.01.2006")
	},
}).ParseFiles("./templates/layout.html.tmpl", "./templates/invoice.html.tmpl"))

func generateFilename() string {
	var bytes = make([]byte, 4)
	rand.Read(bytes[:])
	return fmt.Sprintf("%x", bytes)
}

var chromeFlags = []string{
	"--headless",
	"--disable-accelerated-2d-canvas",
	"--disable-gpu",
	"--allow-pre-commit-input",
	"--disable-background-networking",
	"--disable-background-timer-throttling",
	"--disable-backgrounding-occluded-windows",
	"--disable-breakpad",
	"--disable-client-side-phishing-detection",
	"--disable-component-extensions-with-background-pages",
	"--disable-component-update",
	"--disable-default-apps",
	"--disable-extensions",
	"--disable-features=Translate,BackForwardCache,AcceptCHFrame,MediaRouter,OptimizationHints",
	"--disable-hang-monitor",
	"--disable-ipc-flooding-protection",
	"--disable-popup-blocking",
	"--disable-prompt-on-repost",
	"--disable-renderer-backgrounding",
	"--disable-sync",
	"--enable-automation",
	"--enable-features=NetworkServiceInProcess2",
	"--export-tagged-pdf",
	"--force-color-profile=srgb",
	"--hide-scrollbars",
	"--metrics-recording-only",
	"--no-default-browser-check",
	"--no-first-run",
	"--no-service-autorun",
	"--password-store=basic",
	"--use-mock-keychain",
	"--no-sandbox",
}

const ChromiumExecutable = "/usr/lib/chromium/chromium"

func generatePDFFromSource(source []byte) ([]byte, error) {
	outDir, err := filepath.Abs("./out")
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outDir, 0o777); err != nil {
		return nil, err
	}

	id := generateFilename()
	tmpFile := path.Join(outDir, id+".html")
	outFile := path.Join(outDir, id+".pdf")
	args := append(chromeFlags, fmt.Sprintf(`--print-to-pdf=%s`, outFile), tmpFile)

	os.WriteFile(tmpFile, source, 0o644)

	out := bytes.NewBuffer([]byte{})

	cmd := exec.Command(ChromiumExecutable, args...)
	cmd.Stderr = out
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		log.Print(err)
		return nil, errors.New(out.String())
	}

	return os.ReadFile(outFile)
}

func main() {
	source := bytes.NewBuffer([]byte{})
	err := pdfTemplate.Execute(source, invoiceAssigns{
		Styles:  styles,
		Invoice: generateInvoice(time.Now()),
	})
	if err != nil {
		log.Print(err)
	}

	// Serve the newly generated PDF file to the browser to view the generated PDF
	http.Handle("/pdf", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pdf, err := generatePDFFromSource(source.Bytes())
		if err != nil {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(pdf)
	}))

	http.Handle("/html", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(source.Bytes())
	}))

	http.Handle("/source", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.Write(source.Bytes())
	}))

	log.Print("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
