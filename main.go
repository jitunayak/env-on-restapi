package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"env-on-restapi/colors"
	"env-on-restapi/constants"
	"flag"
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
	shouldStartServer := flag.Bool("server", false, "starts the server")
	cron := flag.Bool("cron", false, "runs cron job")
	interval := flag.Int("interval", 10, "interval for cron job")
	command := flag.String("cmd", "echo no commands passed to run", "command to run periodically")
	portNumber := flag.String("port", "8088", "server port")
	flag.Parse()
	port := fmt.Sprintf(":%s", *portNumber)

	if *shouldStartServer {
		fmt.Println("starting web server")
		fmt.Printf(colors.Red+"Starting server at port %v 🔥 \n\n", port)
		fmt.Printf(colors.Yellow+"GET - http://localhost%v/aws"+colors.Reset, port)
		fmt.Println(colors.Green + constants.Sample_code + colors.Reset)
		startWebServer(port)

	} else {
		fmt.Println("You are on command line. Use eli --help to know all parameters")
		if *cron {
			s := gocron.NewScheduler(time.UTC)
			startCronJobInShell(s, *command, *interval)
			s.StartBlocking()
		}
	}
}

func startWebServer(port string) {

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
			s := gocron.NewScheduler(time.UTC)
			startCronJobInShell(s, command, intervalNumber)
			s.StartAsync()
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

	if err := http.ListenAndServe(port, nil); err != nil {
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

func startCronJobInShell(s *gocron.Scheduler, command string, interval int) {
	fmt.Println(colors.Green+"interval: ", interval, "seconds")
	fmt.Println("command: ", command+colors.Yellow)
	if s.IsRunning() {
		s.Stop()
	}
	fmt.Println("\nstarted corn job" + colors.Reset)

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

	})
}
