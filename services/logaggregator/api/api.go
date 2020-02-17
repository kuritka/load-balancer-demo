package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"lb/common/entity"
	"net/http"
)

/*
because log aggregator is called from any other service, we should provide one API for all other services to not
repeat the code
*/

var (
	//ignoring developer certificates we are using
	//allows us working with HTTPS without any problems with certificates
	transport = http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client = &http.Client{Transport: &transport}

	logServiceUrl = flag.String("logservice", "https://127.0.0.1:6080", "address of logging service")
)

func WriteEntry(entry *entity.LogEntry) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(entry)
	req, _ := http.NewRequest(http.MethodPost, *logServiceUrl, &buf)
	client.Do(req)
}
