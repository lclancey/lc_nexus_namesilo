package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func main() {
	var (
		logStr    string
		logger    *log.Logger
		tempBytes []byte
		config    config
		err       error
	)
	// 0.1 Prepare config
	tempBytes, err = os.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("Read config file err %+v", err)
	}
	err = json.Unmarshal(tempBytes, &config)
	if err != nil {
		log.Fatalf("Parse config file err %+v", err)
	}
	// 0.2 Prepare logger.
	logFile, err := os.OpenFile("./server.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Create log file err %+v", err)
	}
	logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
	// 1. Request from client(typically from a router.)
	//		query should be like: http://THISSERVER:5066/namesiloapi?prefix=aaaa:bbbb:cccc:dddd
	http.HandleFunc("/namesiloapi", func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Recieve a namesilo request.")
		// 2. Get routerPrefix from client query string
		listBody := list_body{}
		var resp *http.Response
		routerPrefix := r.URL.Query().Get("prefix")
		if !regexp.MustCompile(`^[0-9a-fA-F:]+$`).MatchString(routerPrefix) {
			logStr = "Request from client error: prefix not match."
			goto onStuck
		}
		// 3. Request from namesilo for all records.
		resp, err = http.Get(NS_BaseUrl + NS_DnsListPath + "?" + url.Values{
			"version": {"1"},
			"type":    {"xml"},
			"domain":  {config.NS_Domain},
			"key":     {config.NS_Key},
		}.Encode())
		if err != nil {
			logStr = "Http get list err:" + err.Error()
			goto onStuck
		}
		defer resp.Body.Close()
		// 4. Parse XML body to local variable.
		tempBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			logStr = "ListAPI body read err:" + err.Error()
			goto onStuck
		}
		if xml.Unmarshal(tempBytes, &listBody) != nil {
			logStr = "ListAPI body parse err:" + err.Error()
			goto onStuck
		}
		if listBody.Code != 300 || len(listBody.Resource_record) == 0 {
			logStr = "ListAPI reported err code:" + listBody.Detail
			goto onStuck
		}
		// 5. Select useful records to a new variable
		logStr = "Check:"
		for _, localHost := range config.Hosts {
			for _, remoteRecord := range listBody.Resource_record {
				if remoteRecord.FullName == localHost.Name+"."+config.NS_Domain {
					logStr += localHost.Name
					if routerPrefix+localHost.IpSuffix == remoteRecord.FullIP {
						logStr += " SAME "
					} else {
						http.Get(NS_BaseUrl + NS_DnsUpdatePath + "?" + url.Values{
							"version": {"1"},
							"type":    {"xml"},
							"key":     {config.NS_Key},
							"domain":  {config.NS_Domain},
							"rrid":    {remoteRecord.Record_id},
							"rrhost":  {localHost.Name},
							"rrvalue": {routerPrefix + localHost.IpSuffix},
							"rrttl":   {remoteRecord.Ttl},
						}.Encode())
						// fmt.Println(NS_BaseUrl + NS_DnsUpdatePath + "?" + url.Values{
						// 	"version": {"1"},
						// 	"type":    {"xml"},
						// 	"key":     {config.NS_Key},
						// 	"domain":  {config.NS_Domain},
						// 	"rrid":    {remoteRecord.Record_id},
						// 	"rrhost":  {localHost.Name},
						// 	"rrvalue": {routerPrefix + localHost.IpSuffix},
						// 	"rrttl":   {remoteRecord.Ttl},
						// }.Encode())
						logStr += " CHANGE "
					}
				}
			}
		}
		logger.Println(logStr)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, logStr)
		logger.Println("------------------------------------------")
		return
	onStuck:
		logger.Println(logStr)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, logStr)
		logger.Println("------------------------------------------")
	})
	http.ListenAndServe(":5066", nil)
}
