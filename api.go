package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TypeCount struct {
	ID      string  `bson:"_id"`
	Count   int     `bson:"Count"`
	Average float64 `bson:"Average"`
}

var collection *mongo.Collection

func startAPI() {
	/* Router Init */
	router := mux.NewRouter()

	/*Connecting to Database*/
	clientOPtions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOPtions)

	if err != nil {
		log.Fatal(err)
	}

	/* Checking Connection*/
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully Connected to Database")

	collection = client.Database("NYC_DATA").Collection("Buildings")

	/* Endpoints */
	router.HandleFunc("/getBuildings", getBuildings).Methods("GET")
	router.HandleFunc("/updateBuildings", updateBuilding).Methods("PUT")
	router.HandleFunc("/addBuilding", addBuilding).Methods("POST")
	router.HandleFunc("/statHeightByType", statHeightByType).Methods("GET")

	log.Fatal(http.ListenAndServe(":4018", router))
}

func getBuildings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var buildings []Building_Insertable
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var building Building_Insertable
		cursor.Decode(&building)
		buildings = append(buildings, building)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(buildings)

}

func updateBuilding(w http.ResponseWriter, r *http.Request) {

}

func addBuilding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var building Building_Insertable
	_ = json.NewDecoder(r.Body).Decode(&building)
	res, _ := collection.InsertOne(context.TODO(), building)
	json.NewEncoder(w).Encode(res)
}

func statHeightByType(w http.ResponseWriter, r *http.Request) {
	var buildtype []TypeCount
	pipeline := []bson.M{bson.M{"$group": bson.M{"_id": "$type", "Count": bson.M{"$sum": 1}, "Average": bson.M{"$avg": "$height"}}}}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var typecount TypeCount
		cursor.Decode(&typecount)
		buildtype = append(buildtype, typecount)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(buildtype)
}

func averageHeight(w http.ResponseWriter, r *http.Request) {

}
