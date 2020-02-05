package main

import (
	"io"
	"net/http"
	"net/url"
)

var (
	// url of application servers
	appServers []string

	// which application server we just called
	appServerIndex = 0

	// callback to appserver
	client = http.Client{Transport: &transport}
)

func processRequests() {
	for {
		select {
		case request := <-requestCh:
			println("request")
			if len(appServers) == 0 {
				request.w.WriteHeader(http.StatusInternalServerError)
				_, _ = request.w.Write([]byte("No app servers found"))
				request.doneCh <- struct{}{}
				continue
			}
			//todo: refactor
			incrementRoundRobin()
			host := appServers[appServerIndex]

			//handling any single request we have, so function must be as fast as possible
			//thats why in goroutine
			go processRequest(host, request)

		case host := <-registerCh:
			println("register: " + host)
			isFound := false
			for _, h := range appServers {
				if host == h {
					isFound = true
					break
				}
			}
			if !isFound {
				appServers = append(appServers, host)
			}

		case host := <-unregisterCh:
			println("register: " + host)
			for i := len(appServers) - 1; i >= 0; i-- {
				//removing element on the specified position.
				if appServers[i] == host {
					appServers = append(appServers[:i], appServers[i+1:]...)
					break
				}
			}

		case <-heartBeatCh:
			println("heartbeat")
			//copying slice
			servers := appServers[:]
			go func(servers []string) {
				for _, server := range servers {
					resp, err := http.Get("https://" + server + "/ping")
					if err != nil || resp.StatusCode != http.StatusOK {
						//unregister death server
						unregisterCh <- server
					}
				}
			}(servers)
		}
	}
}

//we must build new request from original one, send it to right host and forward back to the requester
func processRequest(host string, request *webRequest) {
	//build url for new host
	hostUrl, _ := url.Parse(request.r.URL.String())
	hostUrl.Scheme = "https"
	hostUrl.Host = host
	println(host)
	println(hostUrl.String())
	req, _ := http.NewRequest(request.r.Method, hostUrl.String(), request.r.Body)
	//because request headers in go is map of slice of strings we must translate into string of headers to new request
	for k, v := range request.r.Header {
		values := ""
		for _, headerValue := range v {
			values += headerValue + " "
		}
		//to slice of strings
		req.Header.Add(k, values)
	}

	resp, err := client.Do(req)
	if err != nil {
		request.w.WriteHeader(http.StatusInternalServerError)
		request.doneCh <- struct{}{}
		return
	}
	//now we have response headers to work with
	//in production we will need some exceptions here
	//we don't want send any headers making security issues within organisation
	for key, header := range resp.Header {
		headers := ""
		for _, headerValue := range header {
			headers += headerValue + " "
		}
		request.w.Header().Add(key, headers)
	}
	_, _ = io.Copy(request.w, resp.Body)

	request.doneCh <- struct{}{}
}

func incrementRoundRobin() {
	//round robin
	appServerIndex++

	if appServerIndex == len(appServers) {
		//restarting round robin
		appServerIndex = 0
	}
}
