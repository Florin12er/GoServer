package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "hello")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve values from form
	name := r.FormValue("name")
	address := r.FormValue("address")

	// Check if name or address is empty
	if name == "" || address == "" {
		http.Error(w, "Name and Address are required fields", http.StatusBadRequest)
		return
	}

	// Connect to MongoDB
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		http.Error(w, "MongoDB URI not found", http.StatusInternalServerError)
		return
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB", http.StatusInternalServerError)
		log.Printf("Failed to connect to MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(context.Background())

	// Insert data into MongoDB
	coll := client.Database("your_database_name").Collection("users") // Adjust your database and collection name
	user := bson.D{
		{"name", name},
		{"address", address},
	}

	_, err = coll.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Failed to insert data into MongoDB", http.StatusInternalServerError)
		log.Printf("Failed to insert data into MongoDB: %v\n", err)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Data inserted into MongoDB:\nName: %s\nAddress: %s\n", name, address)
}
func testHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/test" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, "./static/test.html")
}

func main() {
	// Serve static files from the "./static" directory
	var fileServer = http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	// Define HTTP handlers for specific routes
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/test", testHandler)

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to MongoDB
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " +
			"www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// Ensure MongoDB connection is established before starting HTTP server
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	// Run MongoDB function
	if err := MongoDB(); err != nil {
		log.Fatal("Error running MongoDB function:", err)
	}

	// Start the HTTP server on port 8080
	fmt.Print("The server has started at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func MongoDB() error {
	// Connect to MongoDB
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return fmt.Errorf("MONGODB_URI environment variable is not set")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// Query MongoDB example
	coll := client.Database("sample_mflix").Collection("movies")
	title := "Back to the Future"
	var result bson.M

	// Try to find the document with the specified title
	err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		// Document doesn't exist, insert a new document
		newMovie := bson.D{
			{"title", title},
			{"year", 1985},
			{"genre", "Science Fiction"},
		}
		_, err := coll.InsertOne(context.TODO(), newMovie)
		if err != nil {
			return fmt.Errorf("Error inserting document: %v", err)
		}
		fmt.Println("Inserted new document:", newMovie)
		return nil
	} else if err != nil {
		return fmt.Errorf("Error fetching document: %v", err)
	}

	// Document exists, print JSON result
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return fmt.Errorf("Error encoding JSON: %v", err)
	}
	fmt.Printf("%s\n", jsonData)
	return nil
}

