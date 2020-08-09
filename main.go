package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
)

type PostWriteSchema struct {
	Body string `json:"body"`
}

type IndexResponse struct {
	Id   string `json:"id" bson:"_id"`
	Body string `json:"body"`
}

func main() {
	ctx := context.Background()
	mongoClient := getMongoClient(ctx)

	collection := mongoClient.Database("gblog").Collection("post")

	r := mux.NewRouter()
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			log.Println("could not get all posts")
			writer.WriteHeader(500)
			return
		}

		results := make([]*IndexResponse, 0)
		if err = cursor.All(ctx, &results); err != nil {
			writer.WriteHeader(500)
			return
		}

		jsonResp, err := json.Marshal(results)
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		_, err = writer.Write(jsonResp)
		if err != nil {
			log.Printf("could not write response: %s", request.RequestURI)
		}
		return
	})

	r.HandleFunc("/add/", func(writer http.ResponseWriter, request *http.Request) {
		doc := PostWriteSchema{}

		requestBody, err := ioutil.ReadAll(request.Body)
		err = json.Unmarshal(requestBody, &doc)
		if err != nil {
			return
		}
		_, err = collection.InsertOne(ctx, doc)
		if err != nil {
			return
		}
		writer.WriteHeader(http.StatusCreated)
	})

	r.HandleFunc("/{postId}/", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		postId, ok := vars["postId"]
		if !ok {
			writer.WriteHeader(http.StatusInternalServerError)
		}

		docID, err := primitive.ObjectIDFromHex(postId)
		singleResult := collection.FindOne(ctx, bson.M{"_id": docID})
		if singleResult.Err() != nil {
			switch singleResult.Err() {
			case mongo.ErrNoDocuments:
				writer.WriteHeader(http.StatusNotFound)
			default:
				writer.WriteHeader(http.StatusInternalServerError)
			}
			log.Printf("could not get post: %s", postId)
			return
		}

		results := IndexResponse{}
		if err := singleResult.Decode(&results); err != nil {
			log.Printf("could not decode mongodb response: %s", postId)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonResp, err := json.Marshal(results)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(jsonResp)
		if err != nil {
			log.Printf("could not write response: %s", request.RequestURI)
		}
		return
	})

	http.Handle("/", r)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to listen")
	}
}

func getMongoClient(ctx context.Context) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("could not create mongo client")
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("could not connect mongo client")
	}
	return client
}
