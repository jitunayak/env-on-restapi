package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"env-on-restapi/constants"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
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

	fmt.Println(len(os.Args), os.Args)

	defer color.Unset()

	if *shouldStartServer {
		red := color.New(color.FgBlue).Add(color.Bold).Add(color.BgYellow)

		red.Printf("\n 🦄 starting blazing fast web server on port %v \n\n", port)
		// color.Red("server started at port %v 🔥 \n\n", port)
		color.Yellow("GET - http://localhost%v/aws\n\n", port)
		color.Black(constants.Title)
		color.Green(constants.Sample_code)
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

func getCurrentShell() string {
	switch runtime.GOOS {
	case "windows":
		return "powershell"
	case "darwin":
		return "zsh"
	case "linux":
		return "bash"
	default:
		log.Fatal("no shell found to execute command")
		return ""
	}
}

func startCronJobInShell(s *gocron.Scheduler, command string, interval int) {
	color.Green("interval: ", interval, "seconds")
	color.Green("command: ", command)
	if s.IsRunning() {
		s.Stop()
	}
	fmt.Println("\nstarted corn job")

	s.Every(interval).Seconds().Do(func() {
		cmd := exec.Command(getCurrentShell(), "-c", command)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

	})
}

// ! TODO
func getEliConfigurationPath() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join("", userHomeDir, ".eli", "configuration")
}

// ! TODO
func readConfiguration() {

	if _, err := os.Stat(getEliConfigurationPath()); err == nil {
		os.ReadFile(getEliConfigurationPath())
	}
}

// ! TODO
func updateConfiguration(config string) {

	if _, err := os.Stat(getEliConfigurationPath()); err == nil {

	} else if errors.Is(err, os.ErrNotExist) {
		// os.Mkdir(filepath.Join(userHomeDir, ".eli"), os.ModePerm)
		f, err := os.Create(getAwsCredentialFilePath())
		if err != nil {
			log.Fatal(err)
		}
		f.WriteString(config)
		defer f.Close()
	}
}
