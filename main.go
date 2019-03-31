package main

import (
	"log"
	// "github.com/gorilla/mux"
)

func main() {

	log.Println("Starting ETL Section")

	var buildings []Building
	/* Getting Data From the Open Data API */
	extract(&buildings)
	/* Transforming and Loading into a local instance of MongoDB*/
	trnsfrmLd(buildings)

	log.Println("Successfully Added Entries to Database. ETL End")

	/* REST API */

	log.Println("Starting Rest-API")

	/* Router Init */
	// router := mux.NewRouter()

	/* Endpoints */
	// router.HandleFunc("/buildings", getBuildings).Method("GET")

}
