package main

import (
  "encoding/json"
  "database/sql"
  "fmt"
  "log"
  "net/http"
  "os"

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
  initDb()
  defer db.Close()
  http.HandleFunc("/api/v1/foods", foodsHandler)
  http.HandleFunc("/api/v1/food/{id}", foodHandler)
  log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func foodsHandler(w http.ResponseWriter, r *http.Request) {
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
    //...
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

// package main

// import (
//     "os"
    // "database/sql"
    // "fmt"
    // _ "github.com/lib/pq"
    // "encoding/json"
    // "log"
    // "net/http"
    // "github.com/gorilla/mux"
// )
//

// type Food struct {
//     ID        string   `json:"id,omitempty"`
//     Name      string   `json:"name,omitempty"`
//     Calories  int      `json:"calories,omitempty"`
//  }

// var foods []Food

// func GetFoods(w http.ResponseWriter, r *http.Request) {
//   json.NewEncoder(w).Encode(foods)
// }
//
// func (a *App) getFood(w http.ResponseWriter, r *http.Request) {
//   vars := mux.Vars(r)
//   id, err := strconv.Atoi(vars["id"])
//   if err != nil {
//       respondWithError(w, http.StatusBadRequest, "Invalid product ID")
//       return
//   }
//
//   p := food{ID: id}
//   if err := p.getFood(a.DB); err != nil {
//       switch err {
//       case sql.ErrNoRows:
//           respondWithError(w, http.StatusNotFound, "Food not found")
//       default:
//           respondWithError(w, http.StatusInternalServerError, err.Error())
//       }
//       return
//   }
//
//   respondWithJSON(w, http.StatusOK, p)
// }

// func (a *App) createFood(w http.ResponseWriter, r *http.Request) {
//   // params := mux.Vars(r)
//   var p food
//   decoder := json.NewDecoder(r.Body)
//   if err := decoder.Decode(&p); err != nil {
//       respondWithError(w, http.StatusBadRequest, "Invalid request payload")
//       return
//   }
//   defer r.Body.Close()
//
//   if err := p.createFood(a.DB); err != nil {
//       respondWithError(w, http.StatusInternalServerError, err.Error())
//       return
//   }
//
//   respondWithJSON(w, http.StatusCreated, p)
// }

// func UpdateFood(w http.ResponseWriter, r *http.Request) {}
// func DeleteFood(w http.ResponseWriter, r *http.Request) {}
// func GetMeals(w http.ResponseWriter, r *http.Request) {}
// func GetMealFoods(w http.ResponseWriter, r *http.Request) {}
// func CreateMealFood(w http.ResponseWriter, r *http.Request) {}
// func DeleteMealFood(w http.ResponseWriter, r *http.Request) {}

// Database creation
// CREATE TABLE foods
//     (
//         id serial NOT NULL,
//         name character varying(100) NOT NULL,
//         calories int,
//     );

// our main function
// func main() {
//   a := app{}
//   a.Initalize(
//     os.Getenv("APP_DB_USERNAME"),
//     os.Getenv("APP_DB_PASSWORD"),
//     os.Getenv("APP_DB_NAME"))
//   a.Run(":8080")
  // dbinfo := fmt.Sprintf("user=%s "+
  //   "password=%s dbname=%s sslmode=disable",
  //   DB_USER, DB_PASSWORD, DB_NAME)
  // db, err := sql.Open("postgres", dbinfo)
  // if err != nil {
  //   panic(err)
  // }
  // defer db.Close()
  //
  // err = db.Ping()
  // if err != nil {
  //   panic(err)
  // }
  //
  // fmt.Println("Successfully connected!")
  //
  // router := mux.NewRouter()
  // router.HandleFunc("/api/v1/foods", GetFoods).Methods("GET")
  // router.HandleFunc("/api/v1/foods/{id}", GetFood).Methods("GET")
  // router.HandleFunc("/api/v1/foods", CreateFood).Methods("POST")
  // router.HandleFunc("/api/v1/foods/{id}", UpdateFood).Methods("PATCH")
  // router.HandleFunc("/api/v1/foods/{id}", DeleteFood).Methods("DELETE")
  // router.HandleFunc("/api/v1/meals", GetMeals).Methods("GET")
  // router.HandleFunc("/api/v1/meals/{id}/foods", GetMealFoods).Methods("GET")
  // router.HandleFunc("/api/v1/meals/{id}/foods/{id}", CreateMealFood).Methods("POST")
  // router.HandleFunc("/api/v1/meals/{id}/foods/{id}", DeleteMealFood).Methods("DELETE")
//   log.Fatal(http.ListenAndServe(":8000", router))
// }
