package main

import (
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "log"
  "net/http"
)

// Container is what we use to store values from JSON request to create a new
// container.
type Container struct {
  DbKeyId int
  DbUserId int
}

func main() {
  rtr := mux.NewRouter()
  rtr.HandleFunc("/", Welcome).Methods("GET")
  rtr.HandleFunc("/new", CreateContainer).Methods("POST")
  rtr.HandleFunc("/delete", DeleteContainer).Methods("DELETE")

  http.Handle("/", rtr)

  log.Println("Listening...")
  http.ListenAndServe(":3000", nil)
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
      fmt.Println(newContainer)
    } else {
      fmt.Println("unable to unmarshall the JSON", jsonInputReturn)
    }
  }
}

func DeleteContainer(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Not implemented yet..."))
}
