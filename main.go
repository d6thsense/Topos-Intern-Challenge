package main

import (
	"log"
)

func main() {

	log.Println("Starting ETL Section")

	var buildings []Building
	/* Getting Data From the Open Data API */
	extract(&buildings)
	/* Transforming and Loading into a local instance of MongoDB*/
	trnsfrmLd(buildings)

	log.Println("Successfully Added Entries to Database.")
	log.Println("Ending ETL Section")

	/* REST API */

	log.Println("Starting Rest-API")
	startAPI()
	log.Println("Quiting Rest API")
}
