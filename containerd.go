// Container Daemon - Used to facilitate container creation.
// Copyright (c) 2014, Massachusetts Institute of Technology
// Please see LICENSE.md for licensing information.

package main

import (
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/samalba/dockerclient"
  "log"
  "time"
  "net/http"
  "os"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

// Container is what we use to store values from JSON request to create a new
// container.
type Container struct {
  DbKeyId int
  DbUserId int
  SshKey string
  SshUser string
  SshPort string
  DockerId string
}

// So we don't have to store database config in the host
type Config struct {
  DataHost string
  DataPort int
  DataUser string
  DataPass string
  DataBase string
}

var config Config

// Handling of Docker event callbacks
func eventCallback(event *dockerclient.Event, args ...interface{}) {
    log.Printf("Received event: %#v\n", *event)
}

func main() {
  // Ensure the config file is there.
  configfile := os.Args[1]
  if len(configfile) == 0 {
    fmt.Printf("A configuration file was not specified as the first argument.\n")
    os.Exit(1)
  }

  if _, err := os.Stat(configfile); os.IsNotExist(err) {
    fmt.Printf("Configuration file does not exist: %s\n", configfile)
    os.Exit(1)
  }

  file, _ := os.Open(os.Args[1])
  decoder := json.NewDecoder(file)
  config = Config{}
  err := decoder.Decode(&config)
  if err != nil {
    fmt.Println("Error decoding Config JSON: ", err)
  }

  // Init Gorilla MUX
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
      timerSetup := time.Now()
      timerTotal := time.Now()

      // New Container
      log.Printf("Container is ?", newContainer)

      // Open SQL database

      database := "%s:%s@tcp(%s:%d)/%s"
      dataString := fmt.Sprintf(database, config.DataUser, 
        config.DataPass, config.DataHost, config.DataPort, config.DataBase)

      db, err := sql.Open("mysql", dataString)
      if err != nil {
        log.Fatal(err)
      }
      log.Print("Connected DB")
      
      // Query for Key
      log.Printf("Querying for Key ID ", newContainer.DbKeyId)
      keyQuery := fmt.Sprintf("SELECT sshkey FROM %s.keys WHERE id=?", config.DataBase)
      row := db.QueryRow(keyQuery, newContainer.DbKeyId)
      err = row.Scan(&newContainer.SshKey)
      if err != nil {
        log.Fatal(err)
      }
      log.Printf("Key is ?", newContainer.SshKey)

      // Query for Username
      log.Printf("Querying for Username ", newContainer.DbUserId)
      userQuery := fmt.Sprintf("SELECT login FROM %s.users WHERE id=?", config.DataBase)
      row = db.QueryRow(userQuery, newContainer.DbUserId)
      err = row.Scan(&newContainer.SshUser)
      if err != nil {
        log.Fatal(err)
      }
      log.Printf("User is ?", newContainer.SshUser)
      log.Printf("Setup complete, took ?", time.Since(timerSetup))
      timerDocker := time.Now()
      // We have the details we need to create a container now.
      cmdvars := []string{newContainer.SshUser, newContainer.SshKey}
      containerConfig := &dockerclient.ContainerConfig{
        Image: "sshsession",
        Cmd: cmdvars,
        //ExposedPorts: ports,
      }
      startConfig := &dockerclient.StartConfig{PublishAllPorts: true}

      // Connect to internal socket, create the container
      docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock")
      newContainer.DockerId, err = docker.CreateContainer(containerConfig)
      if err != nil {
        log.Fatal(err)
      }
      log.Printf("Container ID is ?", newContainer.DockerId)

      // Start the actual container
      err = docker.StartContainer(newContainer.DockerId, startConfig)
      if err != nil {
        log.Fatal(err)
      }
      log.Printf("Docker Init complete, took ?", time.Since(timerDocker))
      log.Printf("Docker post and return to webapp start")
      timerReturn := time.Now()
      // Get the hostport info
      containerInfo, err := docker.InspectContainer(newContainer.DockerId)
      if err != nil {
        log.Fatal(err)
      }
      newContainer.SshPort = containerInfo.NetworkSettings.Ports["22/tcp"][0].HostPort

      // Return the docker container info
      webOutput, err := json.Marshal(newContainer)
      fmt.Fprintf(w, string(webOutput))
      fmt.Println(string(webOutput))
      log.Printf("Return to webapp took ?, total process was ?", time.Since(timerReturn), time.Since(timerTotal))

    } else {
      fmt.Println("unable to unmarshall the JSON", jsonInputReturn)
    }
  }
}

func DeleteContainer(w http.ResponseWriter, r *http.Request) {
  p := make([]byte, r.ContentLength)
  _, httpReadReturn := r.Body.Read(p)
  if httpReadReturn == nil {
    timerDelTotal := time.Now()
    timerDelete := time.Now()
    // Create a variable for our new container info, unmarshall
    var stopContainer Container
    jsonInputReturn := json.Unmarshal(p, &stopContainer)
    if jsonInputReturn == nil {
      docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock")
      log.Println("Killing ?", stopContainer.DockerId)
      err := docker.KillContainer(stopContainer.DockerId)
      log.Printf("Delete took ?", time.Since(timerDelete))
      if err != nil {
        log.Fatal(err)
      }
      err = docker.RemoveContainer(stopContainer.DockerId)
      if err != nil {
        log.Fatal(err)
      }
      fmt.Fprintf(w, "{\"result\":\"success\"}")
      log.Printf("Delete Process took ?", time.Since(timerDelTotal))
    }
  }
}
