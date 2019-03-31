package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Building struct {
	BOROUGH     string `json:"bin"`
	LSTSTATTYPE string `json:"lststatype"`
	CONSTRCTYR  string `json:"cnstrct_yr"`
	HEIGHTROOF  string `json:"heightroof"`
	TYPE        string `json:"feat_code"`
}

type Building_Insertable struct {
	ID               int
	Borough          string
	Status           string
	ConstructionYear int
	Height           float64
	Type             string
}

func main() {

	log.Println("Starting ETL Section")
	var buildings []Building
	//Extracting

	/*Used the API Endpint to get the columns we'll be using here*/
	resp, err := http.Get("https://data.cityofnewyork.us/resource/k8ez-gyqp.json?$select=bin,lststatype,cnstrct_yr,heightroof,feat_code")
	for err != nil {
		log.Fatal(err)
		fmt.Println("Error, getting data ! Trying again")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	err = json.Unmarshal(body, &buildings)
	if err != nil {
		log.Fatal(err)
		panic(err.Error())
	}

	log.Println("Successfully Retrieved Data")

	/* Building Type*/
	featType := make(map[int]string)
	featType[2100] = "Building"
	featType[5100] = "Building Under Construction"
	featType[5110] = "Garage"
	featType[2110] = "Skybridge"
	featType[1001] = "Gas Station Canopy"
	featType[1002] = "Storage Tank"
	featType[1003] = "Placeholder"
	featType[1004] = "Auxilary Structure"
	featType[1005] = "Temporary Structure"

	/*Burough Name*/
	borough := make(map[string]string)
	borough["1"] = "Manhattan"
	borough["2"] = "The Bronx"
	borough["3"] = "Brooklyn"
	borough["4"] = "Queens"
	borough["5"] = "Staten Island"

	// Transforming and Loading
	clientOPtions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOPtions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully Connected to Database")

	insertable := Building_Insertable{}
	collection := client.Database("NYC_DATA").Collection("Buildings")
	for i := range buildings {
		insertable.ID = i + 1
		val, typeErr := strconv.Atoi(buildings[i].TYPE)
		if typeErr != nil {
			panic(typeErr.Error())
		}
		insertable.Type = featType[val] /*Transformed Building ID to Type Name */
		height, heightErr := strconv.ParseFloat(buildings[i].HEIGHTROOF, 64)
		if heightErr != nil {
			panic(heightErr.Error())
		}
		insertable.Height = height
		val, yearErr := strconv.Atoi(buildings[i].CONSTRCTYR)
		if yearErr != nil {
			panic(yearErr.Error())
		}
		insertable.ConstructionYear = val
		insertable.Borough = borough[buildings[i].BOROUGH[0:1]] /*Transformed ID to Boriugh Name*/
		insertable.Status = buildings[i].LSTSTATTYPE
		_, err := collection.InsertOne(context.TODO(), insertable)
		if err != nil {
			log.Println(insertable, " couldn't be inserted, Skipping")
		}
		fmt.Println(insertable)
	}

	log.Println("Successfully Added Entries to Database. ETL End")

	log.Println("Starting Rest-API")

}
