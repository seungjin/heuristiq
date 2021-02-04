package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

//var datastoreClient *datastore.Client

func visit_handler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/visit" {
		http.NotFound(w, r)
		return
	}

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

	_, _, err = client.Collection("visits").Add(
		ctx, map[string]interface{}{
			"timestamp":      time.Now(),
			"request-header": r.Header,
		})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}

	//query := client.Collection("visits").OrderBy("timestamp", firestore.Asc).Limit(10)

	iter := client.Collection("visits").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		log.Println(doc.Data())

	}

	fmt.Println("visit")

}
