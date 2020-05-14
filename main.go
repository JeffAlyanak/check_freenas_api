package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jeffalyanak/check_freenas_api/logger"
	"github.com/jeffalyanak/check_freenas_api/model"
)

func main() {
	logging := true

	logger, err := logger.Get()
	if err != nil {
		fmt.Println(err)
		logging = false
	}

	// Handle cli arguments
	hostip := flag.String("hostip", "", "Host IP")
	username := flag.String("username", "", "Username")
	password := flag.String("password", "", "Password")
	insecure := flag.Bool("skipverifytls", false, "Don't verify TLS certs")
	checktype := flag.String("check", "", "Check to perform. Options are: {alerts,storage}")

	storagewarn := flag.Int64("warn", 80, "Storage used % for warning")
	storagecrit := flag.Int64("crit", 90, "Storage used % for critical")

	flag.Parse()

	if !isPercent(*storagecrit) {
		if logging {
			logger.Println("Critical storage must be between 0-100.")
		}
		fmt.Println("Critical storage must be between 0-100.")
		os.Exit(3)
	}
	if !isPercent(*storagewarn) {
		if logging {
			logger.Println("Storage warning must be between 0-100.")
		}
		fmt.Println("Storage warning must be between 0-100.")
		os.Exit(3)
	}

	if *username == "" {
		if logging {
			logger.Println("No Username provided")
		}
		fmt.Println("No Username provided")
		os.Exit(3)
	}
	if *password == "" {
		if logging {
			logger.Println("No password provided")
		}
		fmt.Println("No password provided")
		os.Exit(3)
	}
	if *hostip == "" {
		if logging {
			logger.Println("No Host IP provided")
		}
		fmt.Println("No Host IP provided")
		os.Exit(3)
	}
	if !(*checktype == "alerts" || *checktype == "storage") {
		if logging {
			logger.Println("Check option missing or invalid. valid options are: {alerts,storage}")
		}
		fmt.Println("Check option missing or invalid. valid options are: {alerts,storage}")
		os.Exit(3)
	}

	if *checktype == "alerts" {
		checkAlerts(*hostip, *username, *password, logger, *insecure, logging)
	} else if *checktype == "storage" {
		checkStorage(*hostip, *username, *password, logger, *insecure, *storagewarn, *storagecrit, logging)
	}

}

func checkAlerts(hostip string, username string, password string, l *logger.Logger, insecure bool, logging bool) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	// Struct for holding data
	var a model.Alerts

	// Build strings for request
	apicall := "https://" + hostip + "/api/v1.0/system/alert/"

	// Build request
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", apicall, nil)
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "application/json")

	// Make Request
	resp, err := client.Do(req)
	if err != nil {
		if logging {
			l.Println("Error:" + err.Error())
		}
		fmt.Println("Error:" + err.Error())
		os.Exit(3)
	}
	defer resp.Body.Close()

	// Marshal json data into struct
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &a); err != nil {
		if logging {
			l.Println(err)
		}
		fmt.Println(err)
		os.Exit(3)
	}

	// Check each alert for those of interest.
	for i := 0; i < a.Meta.TotalCount; i++ {
		if statusNotOkay(a.Objects[i].Level) {
			if logging {
				l.Println(a.Objects[i].Level + " - " + a.Objects[i].Message)
			}
			fmt.Println(a.Objects[i].Level + " - " + a.Objects[i].Message)
			if a.Objects[i].Level == "CRITICAL" {
				os.Exit(2)
			} else if a.Objects[i].Level == "WARNING" {
				os.Exit(1)
			}
		}
	}
	if logging {
		l.Println("OK - No warning or critical Alerts")
	}
	fmt.Println("OK - No warning or critical Alerts")
	os.Exit(0)
}

func checkStorage(hostip string, username string, password string, l *logger.Logger, insecure bool, warn int64, crit int64, logging bool) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	// Struct for holding data
	var s model.Storage

	// Build strings for request
	apicall := "https://" + hostip + "/api/v1.0/storage/volume/"

	// Build request
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", apicall, nil)
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "application/json")

	// Make Request
	resp, err := client.Do(req)
	if err != nil {
		if logging {
			l.Println("Error:" + err.Error())
		}
		fmt.Println("Error:" + err.Error())
		os.Exit(3)
	}
	defer resp.Body.Close()

	// Marshal json data into struct
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &s); err != nil {
		if logging {
			l.Println(err)
		}
		fmt.Println("Unmarshall Error: " + err.Error())
		os.Exit(3)
	}

	for i := 0; i < len(s); i++ {
		if s[i].Status != "HEALTHY" {
			if logging {
				l.Println("CRITICAL - " + s[i].Name + " is " + s[i].Status)
			}
			fmt.Println("CRITICAL - " + s[i].Name + " is " + s[i].Status)
			os.Exit(2)
		} else if stringPercentToInt(s[i].UsedPct) >= warn {
			if logging {
				l.Println("WARNING - " + s[i].UsedPct + " storage used on " + s[i].Name)
			}
			fmt.Println("WARNING - " + s[i].UsedPct + " storage used on " + s[i].Name)
			os.Exit(1)
		} else if stringPercentToInt(s[i].UsedPct) >= crit {
			if logging {
				l.Println("CRITICAL - " + s[i].UsedPct + " storage used on " + s[i].Name)
			}
			fmt.Println("CRITICAL - " + s[i].UsedPct + " storage used on " + s[i].Name)
			os.Exit(2)
		}
	}
	if logging {
		l.Println("OK - all storage HEALTHY and below warning level (" + int64ToStringPercent(warn) + ")")
	}
	fmt.Println("OK - all storage HEALTHY and below warning level (" + int64ToStringPercent(warn) + ")")
	os.Exit(0)
}

func stringPercentToInt(str string) int64 {
	i, _ := strconv.Atoi(strings.TrimSuffix(str, "%"))
	return int64(i)
}

func int64ToStringPercent(i int64) string {
	return (strconv.Itoa(int(i)) + "%")
}

func isPercent(num int64) bool {
	if num < 0 || num > 100 {
		return false
	}
	return true
}

func statusNotOkay(s string) bool {
	if s == "OK" || s == "INFO" {
		return false
	}
	return true
}
