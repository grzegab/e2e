package Application

import (
	"e2e/Interface"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func StartTestWebserver() {
	Interface.InfoMsg("Building web server to test api endpoints.", []any{})

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":9821", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

	//@TODO: can edit config? to get better results and check other settings?
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}
