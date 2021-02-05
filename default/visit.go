package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
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

	var output = ""

	iter := client.Collection("visits").
		OrderBy("timestamp", firestore.Desc).Limit(10).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		a, _ := doc.Data()["request-header"].(map[string]interface{})
		output = output + fmt.Sprintf("%v\t%v\t%v\n",
			doc.Data()["timestamp"], a["X-Forwarded-For"], a["user-agent"],
		)
	}

	fmt.Fprint(w, output)

}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
