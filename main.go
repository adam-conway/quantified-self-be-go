package main

import (
  "encoding/json"
  "database/sql"
  "fmt"
  "log"
  "net/http"
  "os"
  "github.com/gorilla/mux"

  _ "github.com/lib/pq"
)

type foodsIndex struct {
	ID         int
	Name       string
	Calories   int
}

type foods struct {
	Foods []foodsIndex
}

var db *sql.DB

const (
  dbhost = "DBHOST"
  dbport = "DBPORT"
  dbuser = "DBUSER"
  dbpass = "DBPASS"
  dbname = "DBNAME"
)

func main() {
  router := mux.NewRouter()
  initDb()
  defer db.Close()
  router.HandleFunc("/api/v1/foods", getFoods).Methods("GET")
  router.HandleFunc("/api/v1/foods", createFood).Methods("POST")
  http.HandleFunc("/api/v1/food/{id}", foodHandler)
  log.Fatal(http.ListenAndServe("localhost:8000", router))
}

func getFoods(w http.ResponseWriter, r *http.Request) {
  foodindex := foods{}

  err := queryDb(&foodindex)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  out, err := json.Marshal(foodindex)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  fmt.Fprintf(w, string(out))
}

func createFood(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  name := r.FormValue("name")
  calories := r.FormValue("calories")

  sqlStatement := `
  INSERT INTO foods (name, calories)
  VALUES ($1, $2)
  RETURNING id`
  id := 0
  err := db.QueryRow(sqlStatement, name, calories).Scan(&id)
  if err != nil {
    panic(err)
  }
  fmt.Println("New record ID is:", id)
}

func queryDb(foodindex *foods) error {
  rows, err := db.Query(`
    SELECT
        id,
        name,
        calories
    FROM foods`)
  if err != nil {
    return err
  }
  defer rows.Close()
  for rows.Next() {
    food := foodsIndex{}
    err = rows.Scan(
      &food.ID,
      &food.Name,
      &food.Calories,
    )
    if err != nil {
      return err
    }
    foodindex.Foods = append(foodindex.Foods, food)
  }
  err = rows.Err()
  if err != nil {
    return err
  }
  return nil
}

func foodHandler(w http.ResponseWriter, r *http.Request) {
}

func initDb() {
  config := dbConfig()
  var err error
  psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
      "password=%s dbname=%s sslmode=disable",
      config[dbhost], config[dbport],
      config[dbuser], config[dbpass], config[dbname])

    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
      panic(err)
    }
    err = db.Ping()
    if err!= nil {
      panic(err)
    }
    fmt.Println("Successfully connected!")
}

func dbConfig() map[string]string {
  conf := make(map[string]string)
  host, ok := os.LookupEnv(dbhost)
  if !ok {
      panic("DBHOST environment variable required but not set")
  }
  port, ok := os.LookupEnv(dbport)
  if !ok {
      panic("DBPORT environment variable required but not set")
  }
  user, ok := os.LookupEnv(dbuser)
  if !ok {
      panic("DBUSER environment variable required but not set")
  }
  password, ok := os.LookupEnv(dbpass)
  if !ok {
      panic("DBPASS environment variable required but not set")
  }
  name, ok := os.LookupEnv(dbname)
  if !ok {
      panic("DBNAME environment variable required but not set")
  }
  conf[dbhost] = host
  conf[dbport] = port
  conf[dbuser] = user
  conf[dbpass] = password
  conf[dbname] = name
  return conf
}
