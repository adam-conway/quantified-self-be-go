package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    // "strconv"
)

type Food struct {
    ID        string   `json:"id,omitempty"`
    Name      string   `json:"name,omitempty"`
    Calories  int      `json:"calories,omitempty"`
 }

var foods []Food

func GetFoods(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(foods)
}

func GetFood(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
    for _, item := range foods {
    if item.ID == params["id"] {
      json.NewEncoder(w).Encode(item)
      return
    }
  }
  json.NewEncoder(w).Encode(&Food{})
}

func CreateFood(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  var food Food
  _ = json.NewDecoder(r.Body).Decode(&food)
  food.ID = "4"
  food.Name = params["name"]
  // food.Calories = strconv.Atoi(params["Calories"])
  foods = append(foods, food)
  json.NewEncoder(w).Encode(foods)
}

func UpdateFood(w http.ResponseWriter, r *http.Request) {}
func DeleteFood(w http.ResponseWriter, r *http.Request) {}
func GetMeals(w http.ResponseWriter, r *http.Request) {}
func GetMealFoods(w http.ResponseWriter, r *http.Request) {}
func CreateMealFood(w http.ResponseWriter, r *http.Request) {}
func DeleteMealFood(w http.ResponseWriter, r *http.Request) {}

// our main function
func main() {
  router := mux.NewRouter()
  foods = append(foods, Food{ID: "1", Name: "Oreos", Calories: 100})
  foods = append(foods, Food{ID: "2", Name: "Pizza", Calories: 200})
  foods = append(foods, Food{ID: "3", Name: "Taco Bell", Calories: 300})
  router.HandleFunc("/api/v1/foods", GetFoods).Methods("GET")
  router.HandleFunc("/api/v1/foods/{id}", GetFood).Methods("GET")
  router.HandleFunc("/api/v1/foods", CreateFood).Methods("POST")
  // router.HandleFunc("/api/v1/foods/{id}", UpdateFood).Methods("PATCH")
  router.HandleFunc("/api/v1/foods/{id}", DeleteFood).Methods("DELETE")
  // router.HandleFunc("/api/v1/meals", GetMeals).Methods("GET")
  // router.HandleFunc("/api/v1/meals/{id}/foods", GetMealFoods).Methods("GET")
  // router.HandleFunc("/api/v1/meals/{id}/foods/{id}", CreateMealFood).Methods("POST")
  // router.HandleFunc("/api/v1/meals/{id}/foods/{id}", DeleteMealFood).Methods("DELETE")
  log.Fatal(http.ListenAndServe(":8000", router))
}
