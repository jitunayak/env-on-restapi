package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

type AppConfigProperties map[string]string

func main() {
	config := AppConfigProperties{}

	http.HandleFunc("/aws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		shoudlReAuthenticate := r.URL.Query().Get("reAuthenticate")
		interval := r.URL.Query().Get("interval")
		command := r.URL.Query().Get("command")

		if shoudlReAuthenticate == "" {
			shoudlReAuthenticate = "false"
		}
		if interval == "" {
			interval = "3000" //Default Time Out : 50 minutes
		}
		if command == "" && shoudlReAuthenticate == "true" {
			http.Error(w, "command is missing to authenticate", http.StatusBadRequest)
		}

		if shoudlReAuthenticate == "true" {
			intervalNumber, err := strconv.Atoi(interval)
			if err != nil {
				panic("interval time is invaliad")
			}
			startCronJobInShell(command, intervalNumber)
		}

		config := getAwsConfiguration(config)
		data := map[string]interface{}{
			"accessKeyId":  config["aws_access_key_id"],
			"secretKey":    config["aws_secret_access_key"],
			"sessionToken": config["aws_session_token"],
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

func getAwsCredentialFilePath() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	awsCredPath := filepath.Join(userHomeDir, ".aws", "credentials")
	return awsCredPath
}
func getAwsConfiguration(config AppConfigProperties) AppConfigProperties {
	awsCredPath := getAwsCredentialFilePath()
	file, err := os.Open(awsCredPath)

	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}
	return config
}

func startCronJobInShell(command string, interval int) {
	s := gocron.NewScheduler(time.UTC)
	if s.IsRunning() {
		s.Stop()
	}
	fmt.Println("started corn job")

	currentShell := "zsh"

	if runtime.GOOS == "windows" {
		currentShell = "powershell"
	}
	s.Every(interval).Seconds().Do(func() {
		cmd := exec.Command(currentShell, "-c", command)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		// os.Mkdir("jitu", os.ModePerm)

	})
	s.StartAsync()
}
