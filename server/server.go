package main

import (
	"html"
	"html/template"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"os"
	"fmt"
	"time"
)

// Compile templates on start.
var templates = template.Must(template.ParseFiles("static/index.html"))


// The input to the webpage template, carries texada output.
type Output struct {
    OutputTitle string
    OutputResult string
    OutputDisplay string
}

// Where to find the texada binary
var texadaCmd string

// Display the template.
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Outputs an error string.
func printErrorStr(s string) {
	fmt.Println(time.Now(), "Error: " + s)
}

// Outputs the received http request to stdout.
func printRequest(r *http.Request) {
	fmt.Println(time.Now(), "REQ:", html.EscapeString(r.URL.Path), ", ip/port: ", r.RemoteAddr)
}

// Handles POST to /texada/mine/
func mineHandler(w http.ResponseWriter, r *http.Request) {
	printRequest(r)

	// Retrieve form data.
	body := r.FormValue("body")
	args := r.FormValue("args")
	
	// Remove carriage returns from the input, and convert to []byte
	body = strings.Replace(body,"\r","",-1)
	args = strings.Replace(args,"\r","",-1)
	bodybytes := []byte(body)
	argsbytes := []byte(args)

	// Save body to file.
	logfile, err := ioutil.TempFile("/tmp/", "log-");
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ioutil.WriteFile(logfile.Name(), bodybytes, 0600) != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save args to file.
	argsfile, err := ioutil.TempFile("/tmp/", "args-");
	if err != nil {
		os.Remove(logfile.Name())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ioutil.WriteFile(argsfile.Name(), argsbytes, 0600) != nil {
		os.Remove(logfile.Name())
		os.Remove(argsfile.Name())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	outbytes, err := exec.Command(texadaCmd, "-c", argsfile.Name(), logfile.Name()).Output();
	if err != nil {
		os.Remove(logfile.Name())
		os.Remove(argsfile.Name())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	os.Remove(logfile.Name())
	os.Remove(argsfile.Name())
	outstr := strings.TrimSpace(string(outbytes))
	// cmdfull := cmd + " -c " + argsfile.Name() + " " + logfile.Name()

	result := Output{OutputTitle: "Texada output:", OutputResult: outstr, OutputDisplay: "block"}
	renderTemplate(w, "index", result)
}

// Handles GET of /texada/
func viewHandler(w http.ResponseWriter, r *http.Request) {
	printRequest(r)

	result := Output{OutputTitle: "", OutputResult: "", OutputDisplay: "none"}
	renderTemplate(w, "index", result)
}

func main() {
	// Handle command line args:
	args := os.Args
	if len(args) != 3 {
		printErrorStr("Wrong number of arguments.")
		fmt.Println("Usage: go run server.go [texada-bin] [port]")
		fmt.Println("Examples:")
		fmt.Println("\tgo run server.go /Users/ivan/repos/texada/texada 8080")
		fmt.Println("\tgo run server.go /home/bestchai/texada/texada 8080")
		os.Exit(-1)
	}

	texadaCmd = args[1]
	port := args[2]

	fmt.Printf("Using texada binary: %s\n", texadaCmd)
	fmt.Printf("Using port: %s\n", port)

	// Register handlers for different paths.
	http.HandleFunc("/texada/", viewHandler)
	http.HandleFunc("/texada/mine/", mineHandler)
	
	// Register static files handler
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Listen on port with default IP.
	http.ListenAndServe(":" + port, nil)
}