package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type PromResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func main() {

	// pull conifg from envvars
	promURL := os.Getenv("PROMETHEUS_URL")
	if promURL == "" {
		promURL = "http://monitoring-kube-prometheus-prometheus.monitoring.svc.cluster.local:9090"
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	destEmail := os.Getenv("DEST_EMAIL")

	// query prom
	query := `count(kube_node_status_condition{condition="Ready",status="true"})`
	promEndpoint := fmt.Sprintf("%s/api/v1/query?query=%s", promURL, query)

	resp, err := http.Get(promEndpoint)
	if err != nil {
		log.Fatalf("Failed to reach Prometheus: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var promData PromResponse
	if err := json.Unmarshal(body, &promData); err != nil {
		log.Fatalf("Failed to parse Prometheus JSON: %v", err)
	}

	// extract metric value
	nodeCount := "0"
	if len(promData.Data.Result) > 0 && len(promData.Data.Result[0].Value) > 1 {
		nodeCount = fmt.Sprintf("%v", promData.Data.Result[0].Value[1])
	}

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: K3s Alert\r\n"+
		"\r\n"+
		"%s\r\n", destEmail, messageBody))

	// dispatch sms via gmail smtp
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err = smtp.SendMail(smtpAddr, auth, smtpUser, []string{destEmail}, msg)
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}

	log.Println("Successfully transmitted cluster report to phone!")
}
