package main

import (
	"fmt"
	"log"
)

func main() {

	var buildings []Building
	var etlCheck int
	fmt.Println("Press 1 to start ETL process and 0 to skip!")
	_, err := fmt.Scan(&etlCheck)
	if err != nil {
		fmt.Println(err.Error(), "Please re-run")
	}
	if etlCheck != 0 {
		log.Println("Starting ETL Section")
		/* Getting Data From the Open Data API */
		extract(&buildings)
		/* Transforming and Loading into a local instance of MongoDB*/
		trnsfrmLd(buildings)
		log.Println("Ending ETL Section")
	}

	/* REST API */

	log.Println("Starting Rest-API")
	startAPI()
	log.Println("Quiting Rest API")
}
