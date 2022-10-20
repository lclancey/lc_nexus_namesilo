package main

const (
	NS_BaseUrl       = "https://www.namesilo.com/api"
	NS_DnsListPath   = "/dnsListRecords"
	NS_DnsUpdatePath = "/dnsUpdateRecord"
)

type config struct {
	NS_Domain string     `json:"ns_domain"`
	NS_Key    string     `json:"ns_key"`
	Hosts     []cnf_host `json:"hosts"`
}
type cnf_host struct {
	Name     string `json:"name"`
	IpSuffix string `json:"ipsufix"`
}

// jsonFile, err := os.Open("./config.json")
// if err != nil {
// 	fmt.Println(err)
// }
// defer jsonFile.Close()
// var result config
// data, _ := io.ReadAll(jsonFile)
// json.Unmarshal(data, &result)

// record of namesilo API response
type list_record struct {
	Record_id   string `xml:"record_id"`
	Recort_type string `xml:"type"`
	FullName    string `xml:"host"`
	FullIP      string `xml:"value"`
	Ttl         string `xml:"ttl"`
}

type list_body struct {
	Code            uint          `xml:"reply>code"`
	Detail          string        `xml:"reply>detail"`
	Resource_record []list_record `xml:"reply>resource_record"`
}
