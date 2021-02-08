package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"github.com/seungjin/heuristiq/default/fortunecookie"
	"github.com/seungjin/heuristiq/default/visits"
)

func main() {

<<<<<<< HEAD
	mux := http.NewServeMux()

	rootHandler := http.HandlerFunc(root_handler)
	mux.Handle("/", visitLog(rootHandler))

	fortunecookieHandler := http.HandlerFunc(fortunecookie.Fc_handler)
	mux.Handle("/fortunecookie", visitLog(fortunecookieHandler))

	visitsHandler := http.HandlerFunc(visits.Visit_handler)
	mux.Handle("/visits", visitLog(visitsHandler))
=======
	http.HandleFunc("/", root_handler)
	http.HandleFunc("/visit", visit_handler)
	//http.HandleFunc("/fortunecookie", fortunecookie_handler)
>>>>>>> 4a58d6baf1d7e3c3670b4f6ca1954849787a9897

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}

}

func visitLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set this in app.yaml when running in production.
		projectID := os.Getenv("GCLOUD_DATASET_ID")

		// Use the application default credentials
		ctx := context.Background()
		conf := &firebase.Config{ProjectID: projectID}
		app, err := firebase.NewApp(ctx, conf)
		if err != nil {
			log.Fatalln(err)
		}

		client, err := app.Firestore(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		defer client.Close()

		_, _, err = client.Collection("visits-log").Add(
			ctx, map[string]interface{}{
				"timestamp":      time.Now(),
				"method":         r.Method,
				"remote_addr":    r.RemoteAddr,
				"request_uri":    r.RequestURI,
				"host":           r.Host,
				"proto":          r.Proto,
				"request_header": r.Header,
			})
		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		}

		next.ServeHTTP(w, r)
	})
}

func root_handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Fprint(w, "Hello, World!")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		address := r.FormValue("address")
		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Address = %s\n", address)
	default:
		fmt.Fprintf(w, "Not supported HTTP method")
	}
}
