package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("POST /render/{template}", render)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("function not implemented"))
	})

	println("running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func render(w http.ResponseWriter, r *http.Request) {
	template := r.PathValue("template")

	err := os.RemoveAll("documents")
	if err != nil {
		serverError(w, err, "error removing old git path")
		return
	}

	output, err := exec.Command("git", "clone", "--depth=1", "https://github.com/swablab/documents.git").CombinedOutput()
	if err != nil {
		serverError(w, err, "git clone error: "+string(output))
		return
	}

	paramsFile, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(w, err, "error while reading request body")
		return
	}

	err = os.WriteFile("documents/"+template+".yml", paramsFile, 0644)
	if err != nil {
		serverError(w, err, "error while writing param file")
		return
	}

	output, err = exec.Command("typst", "compile", "--root=documents", "documents/"+template+".typ").CombinedOutput()
	if err != nil {
		serverError(w, err, "typst error: "+string(output))
		return
	}

	pdfFile, err := os.ReadFile("documents/" + template + ".pdf")
	if err != nil {
		serverError(w, err, "error while reading pdf")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(pdfFile)
}

func serverError(w http.ResponseWriter, err error, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
	w.Write([]byte(err.Error()))
}
