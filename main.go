package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var todoCollection *mongo.Collection

type Todo struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Content     string             `json:"content"`
	CreateDate  string             `json: "createDate"`
	Deadline    string             `json: "deadline"`
	IsDone      bool               `json: "isDone"`
}
type Todos []Todo

func addTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	todo.CreateDate = toISOString(time.Now())
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := todoCollection.InsertOne(ctx, todo)
	if err != nil {
		log.Printf("Error on creating todo: %v", err)
	}
	json.NewEncoder(w).Encode(result)
}
func allTodos(w http.ResponseWriter, r *http.Request) {
	var todos Todos
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := todoCollection.Find(ctx, bson.M{})
	if err != nil {
		reject(w, err).send()
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}

	if err := cursor.Err(); err != nil {
		reject(w, err).send()
		return
	}
	resolve(w, todos).status(200).send()
}
func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var todo Todo
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := todoCollection.FindOne(ctx, Todo{ID: id}).Decode(&todo)
	if err != nil {
		reject(w, err).send()
		return
	}
}
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
	fmt.Println("GET / endpoint hit")
}
func toISOString(t time.Time) string {
	return t.Format(time.RFC3339)
}
func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/all", allTodos).Methods("GET")
	router.HandleFunc("/todo", addTodo).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}
func initalize() {
	fmt.Println("Listening on port 8081")

	uri := "mongodb://localhost:27017"
	var err error
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Can not connect to db: %v", err)
	}
	todoCollection = client.Database("localdev").Collection("todo")
	handleRequests()
}
func main() {
	initalize()
}
