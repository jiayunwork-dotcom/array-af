package main

import (
	"embed"
	"flag"
	"log"
	"net/http"

	"array-af/internal/web"
)

//go:embed web
var webFS embed.FS

//go:embed example/broadside-8.json
var broadsideJSON []byte

func main() {
	httpAddr := flag.String("http", ":8080", "serve the array-af web console on this address (e.g. :8080)")
	flag.Parse()

	handler := web.NewServer(web.Assets{
		WebFS: webFS,
		Examples: map[string][]byte{
			"broadside-8": broadsideJSON,
		},
	})

	log.Printf("array-af web console on http://localhost%s", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, handler); err != nil {
		log.Fatal(err)
	}
}
