package main
import (
	"net/http"
	"log"
)


func main() {
	fs := http.FileServer(http.Dir("imgs"))
	http.Handle("/", fs)

	log.Println("Listening on 3000...")
	http.ListenAndServe(":3000", nil)
}
