package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {

	type User struct {
		ID     string
		Name   string
		Date   string
		Email  string
		Income string
		Ip     string
	}

	users := []User{}
	jsonFile, err := os.Open("users.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	jsonParser := json.NewDecoder(jsonFile)
	jsonParser.Decode(&users)

	fmt.Println("Successfully Loaded users from json file")

	defer jsonFile.Close()

	http.HandleFunc("/getData", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		switch r.Method {
		case http.MethodGet:
			// write response as formatted json
			json.NewEncoder(w).Encode(users)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/postData", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		switch r.Method {
		case http.MethodPost:

			// add user to users.json
			var user User
			json.NewDecoder(r.Body).Decode(&user)
			user.ID = fmt.Sprintf("%d", len(users))
			users = append(users, user)
			jsonFile, err := os.OpenFile("users.json", os.O_RDWR, 0666)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			defer jsonFile.Close()
			// write new user to users.json
			jsonFile.Seek(0, 0)
			jsonFile.Truncate(0)
			json.NewEncoder(jsonFile).Encode(users)
			fmt.Fprintf(w, "Successfully added user to json file")
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/deleteRow", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		switch r.Method {
		case http.MethodPost:
			// get user id from body
			var user User
			json.NewDecoder(r.Body).Decode(&user)

			// print user to be deleted
			fmt.Println("User to be deleted: ", r.Body)

			// delete user from users.json
			for i, v := range users {
				if v.ID == user.ID {
					users = append(users[:i], users[i+1:]...)
				}
			}
			jsonFile, err := os.OpenFile("users.json", os.O_RDWR, 0666)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			defer jsonFile.Close()
			// write new user to users.json
			jsonFile.Seek(0, 0)
			jsonFile.Truncate(0)
			json.NewEncoder(jsonFile).Encode(users)

			fmt.Fprintf(w, "Successfully deleted user from json file")
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
