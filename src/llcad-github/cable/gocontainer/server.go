package main

import (
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "log"
  "net/http"
  _ "github.com/lib/pq"
  "database/sql"
)

// Container is what we use to store values from JSON request to create a new
// container.
type Container struct {
  DbKeyId int
  DbUserId int
  sshKey string
  sshUser string
  dockerId string
}

func main() {
  rtr := mux.NewRouter()
  rtr.HandleFunc("/", Welcome).Methods("GET")
  rtr.HandleFunc("/new", CreateContainer).Methods("POST")
  rtr.HandleFunc("/delete", DeleteContainer).Methods("DELETE")

  http.Handle("/", rtr)

  log.Println("Listening...")
  http.ListenAndServe(":5000", nil)
}

func Welcome(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello - I am a container creator. Why are you here?\n"))
  w.Write([]byte("I only exist for API access..."))
}

func CreateContainer(w http.ResponseWriter, r *http.Request) {
  // Make a multi-byte array the size of the http content length
  p := make([]byte, r.ContentLength)

  // Read into p the http body
  _, httpReadReturn := r.Body.Read(p)

  if httpReadReturn == nil {
    // Create a variable for our new container info, unmarshall
    var newContainer Container
    jsonInputReturn := json.Unmarshal(p, &newContainer)
    if jsonInputReturn == nil {
      log.Printf("Container is ?", newContainer)

      // we need to get all pgsql'y up in here
      db, err := sql.Open("postgres", 
                          "user=pa20690 dbname=pa20690 sslmode=disable")
      if err != nil {
        log.Fatal(err)
      }
      log.Print("Connected DB")
      log.Printf("Querying for Key ID ", newContainer.DbKeyId)
      row := db.QueryRow("SELECT sshkey FROM keys WHERE id = $1", 
                         newContainer.DbKeyId)
      err = row.Scan(&newContainer.sshKey)
      if err != nil {
        log.Fatal(err)
      }
      log.Printf("Key is ?", newContainer.sshKey)
      log.Printf("Querying for Username ", newContainer.DbUserId)
      row = db.QueryRow("SELECT login FROM users WHERE id = $1",
                         newContainer.DbUserId)
      err = row.Scan(&newContainer.sshUser)
      if err != nil {
        log.Fatal(err)
      }
      log.Printf("User is ?", newContainer.sshUser)
    } else {
      fmt.Println("unable to unmarshall the JSON", jsonInputReturn)
    }
  }
}

func DeleteContainer(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Not implemented yet..."))
}
