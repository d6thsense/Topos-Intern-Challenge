# Topos Intern Challenge

## Dependencies

* Mongo DB Community Edition 4.0.8 [Installation Instructions](www.google.com "Mongo DB")

* Mongo DB Golang Driver 1.0.0
  * [Installation Instructions](https://github.com/mongodb/mongo-go-driver#installation)
  Or
  * You can just run:
  
    ```bash
    go get go.mongodb.org/mongo-driver
    ```
    The output of this may look like a warning stating something like package go.mongodb.org/mongo-driver: no Go files in (...). This is expected .

* Gorilla Mux Router [Installation Instructions](https://github.com/gorilla/mux#install)

## Usage

```bash
go run main.go etl.go api.go
```
The above command will launch the program which will give an option to run the _**ETL process**_ which fetches the data from  Building Footprints Dataset and loads it into MongoDB in a _NYC_DATA.Buildings_ the [0/1 Option]. Then the _**API**_ will be launched which will load the data from the 

### Client
When the process is launched as above, the server will be running on 
[http://localhost:4018]. Thea web client which can run all end-points from the API such as

  * Get All Data
  * Get Stats By Type
  * Get Stats by Borough
  * Get stats by Construction Year
  * Get Global Stats
  * Delete a building &
  * Add a building

The client is made just to test out the API endpoints quickly.