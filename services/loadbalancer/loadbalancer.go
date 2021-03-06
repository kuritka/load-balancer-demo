package loadbalancer

/*
L7 Load balancer implementation.
LB is just simple proxy with added functionality like:

	-	balance traffic
	-	service discovery (to know where to balance)
	-	service termination
	-	heartbeat (to know which service to kill)
*/

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	lbkey     = "/etc/lb/certs/key.pem"
	lbcert    = "/etc/lb/certs/cert.pem"
	discokey  = "/etc/lb/certs/key.pem"
	discocert = "/etc/lb/certs/cert.pem"
)

type webRequest struct {
	r      *http.Request
	w      http.ResponseWriter
	doneCh chan struct{}
}

var (
	requestCh = make(chan *webRequest)

	registerCh = make(chan string)

	unregisterCh = make(chan string)

	//heart beat probe channel
	heartBeatCh = time.Tick(5 * time.Second)
)

var (
	//ignoring developer certificates we are using
	//allows us working with HTTPS without any problems with certificates
	transport = http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
)

func init() {
	http.DefaultClient = &http.Client{Transport: &transport}
}

func Run(lbPort, discoPort string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		doneCh := make(chan struct{})
		requestCh <- &webRequest{r: r, w: w, doneCh: doneCh}
		//waits until LoadBalancer resend request to chosen app server and resend response back
		//or error happens
		<-doneCh
	})

	go processRequests()

	//load balancing
	//if nil DefaultServerMux is used and DSM gets registered handler
	go http.ListenAndServeTLS(lbPort, lbcert, lbkey, nil)

	//service discovery
	http.ListenAndServeTLS(discoPort, discocert, discokey, new(discoHandler))

	//log.Println("server started, press <ENTER> to exit")
	//_, _ = fmt.Scanln()
}

type discoHandler struct{}

func (h *discoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//<ip>:<outgoing port>
	ip := strings.Split(r.RemoteAddr, ":")[0]
	//getting incomming port
	port := r.URL.Query().Get("port")

	switch r.URL.Path {
	case "/register":
		registerCh <- fmt.Sprintf("%v:%v", ip, port)
	case "/unregister":
		unregisterCh <- fmt.Sprintf("%v:%v", ip, port)
	}
}
