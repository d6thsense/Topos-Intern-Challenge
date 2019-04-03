package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TypeCount struct {
	ID                string  `bson:"_id"`
	Count             int     `bson:"Count"`
	AverageHeight     float64 `bson:"avHeight"`
	MinimumHeight     float64 `bson:"minHeight"`
	MaximumHeight     float64 `bson:"maxHeight"`
	StandardDeviation float64 `bson:"stdDev"`
}

type HeightCount struct {
	ID                *Interval `bson:"_id"`
	Count             int       `bson:"count"`
	AverageHeight     float64   `bson:"avHeight"`
	MinimumHeight     float64   `bson:"minHeight"`
	MaximumHeight     float64   `bson:"maxHeight"`
	StandardDeviation float64   `bson:"stdDev"`
}

type Interval struct {
	Minimum int `bson:"min"`
	Maximum int `bson:"max"`
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
	router.HandleFunc("/addBuilding", addBuilding).Methods("POST")
	router.HandleFunc("/removeBuilding/{id:[a-z0-9]*}", removeBuilding).Methods("DELETE")
	router.HandleFunc("/statHeightByType", statHeightByType).Methods("GET")
	router.HandleFunc("/statHeightByYear", statHeightByYear).Methods("GET")
	router.HandleFunc("/statHeightByBorough", statHeightByBorough).Methods("GET")
	router.HandleFunc("/globalStatistics", globalStatistics).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Fatal(http.ListenAndServe(":4018", router))
}

func getBuildings(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving All Buildings")
	w.Header().Set("content-type", "application/json")
	var buildings []Building_Outputable
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var building Building_Outputable
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

/* Add Building*/

func addBuilding(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding a new building")
	w.Header().Set("Content-Type", "application/json")
	var building Building_Insertable
	_ = json.NewDecoder(r.Body).Decode(&building)
	res, _ := collection.InsertOne(context.TODO(), building)
	json.NewEncoder(w).Encode(res)
}

func removeBuilding(w http.ResponseWriter, r *http.Request) {
	log.Println("Removing a building")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	objid, _ := primitive.ObjectIDFromHex(params["id"])
	iddoc := bson.D{{"_id", objid}}
	res, _ := collection.DeleteOne(context.TODO(), iddoc)
	json.NewEncoder(w).Encode(res)
}

/*Statistics by Type*/
func statHeightByType(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Stats by Height")
	var buildtype []TypeCount /* Final Result Holder*/
	pipeline := []bson.M{bson.M{"$group": bson.M{"_id": "$type",
		"Count":     bson.M{"$sum": 1},
		"avHeight":  bson.M{"$avg": "$height"},
		"minHeight": bson.M{"$min": "$height"},
		"maxHeight": bson.M{"$max": "$height"},
		"stdDev":    bson.M{"$stdDevPop": "$height"}}}}
	cursor, err := collection.Aggregate(context.TODO(), pipeline) /*Cursor to iterate over documents returned*/
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(context.TODO())
	/* Iterating over each document */
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
	/*Returning the Final Answer*/
	json.NewEncoder(w).Encode(buildtype)
}

/*Statistics by Boroughs*/
func statHeightByBorough(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Stats by Borough")
	var buildtype []TypeCount /* Final Result Holder*/
	pipeline := []bson.M{bson.M{"$group": bson.M{"_id": "$borough",
		"Count":     bson.M{"$sum": 1},
		"avHeight":  bson.M{"$avg": "$height"},
		"minHeight": bson.M{"$min": "$height"},
		"maxHeight": bson.M{"$max": "$height"},
		"stdDev":    bson.M{"$stdDevPop": "$height"}}}}
	cursor, err := collection.Aggregate(context.TODO(), pipeline) /*Cursor to iterate over documents returned*/
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(context.TODO())
	/* Iterating over each document */
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
	/*Returning the Final Answer*/
	json.NewEncoder(w).Encode(buildtype)
}

/*Statistics by Year Intervals*/
func statHeightByYear(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Stats by Year")
	var buildtype []HeightCount /* Final Result Holder*/
	pipeline := []bson.M{bson.M{
		"$bucketAuto": bson.M{
			"groupBy": "$constructionyear",
			"buckets": 25,
			"output": bson.M{
				"count":     bson.M{"$sum": 1},
				"avHeight":  bson.M{"$avg": "$height"},
				"minHeight": bson.M{"$min": "$height"},
				"maxHeight": bson.M{"$max": "$height"},
				"stdDev":    bson.M{"$stdDevPop": "$height"}}}}}

	cursor, err := collection.Aggregate(context.TODO(), pipeline) /*Cursor to iterate over documents returned*/
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(context.TODO())
	/* Iterating over each document */
	for cursor.Next(context.TODO()) {
		var typecount HeightCount
		cursor.Decode(&typecount)
		buildtype = append(buildtype, typecount)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	/*Returning the Final Answer*/
	json.NewEncoder(w).Encode(buildtype)
}

func globalStatistics(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Global Stats")
	var buildtype []TypeCount /* Final Result Holder*/
	pipeline := []bson.M{bson.M{"$group": bson.M{"_id": "null",
		"Count":     bson.M{"$sum": 1},
		"avHeight":  bson.M{"$avg": "$height"},
		"minHeight": bson.M{"$min": "$height"},
		"maxHeight": bson.M{"$max": "$height"},
		"stdDev":    bson.M{"$stdDevPop": "$height"}}}}
	cursor, err := collection.Aggregate(context.TODO(), pipeline) /*Cursor to iterate over documents returned*/
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(context.TODO())
	/* Iterating over each document */
	for cursor.Next(context.TODO()) {
		var typecount TypeCount
		cursor.Decode(&typecount)
		typecount.ID = "Global"
		buildtype = append(buildtype, typecount)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	/*Returning the Final Answer*/
	json.NewEncoder(w).Encode(buildtype)
}
