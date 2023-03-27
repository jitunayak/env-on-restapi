package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	var Reset = "\033[0m"
	var Red = "\033[31m"
	var Green = "\033[32m"
	var Yellow = "\033[33m"
	// var Blue = "\033[34m"
	// var Purple = "\033[35m"
	var Cyan = "\033[36m"
	// var Gray = "\033[37m"
	// var White = "\033[97m"

	title := `Sample Code Snippet For Postman Test 🦄
-------------------------------------------------------------------`
	aws_url := `GET - http://localhost:8088/aws`
	sample_code := `

const {accessKeyId, secretKey, sessionToken} = pm.response.json();
// for setting global level variables
pm.globals.set("accessKeyId", accessKeyId);
pm.globals.set("secretKey", secretKey);
pm.globals.set("sessionToken", sessionToken);

// or collection level variables
pm.collectionVariables.set("accessKeyId", accessKeyId);
pm.collectionVariables.set("secretKey", secretKey);
pm.collectionVariables.set("sessionToken", sessionToken);

-------------------------------------------------------------------

`
	http.HandleFunc("/aws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
		sessionToken := os.Getenv("AWS_SESSION_TOKEN")

		data := map[string]interface{}{
			"accessKeyId":  accessKeyId,
			"secretKey":    secretKey,
			"sessionToken": sessionToken,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}
		w.Write(jsonData)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		var userRequest = make(map[string]string)
		var data = make(map[string]string)

		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for i, request := range userRequest {
			env := os.Getenv(request)
			data[i] = env
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}
		w.Write(jsonData)
	})

	fmt.Printf(Red + "Starting server at port 8088 🔥 \n\n")
	fmt.Println(Cyan + title + Reset)
	fmt.Println(Yellow + aws_url + Reset)
	fmt.Println(Green + sample_code + Reset)
	if err := http.ListenAndServe(":8088", nil); err != nil {
		log.Fatal(err)
	}
}
