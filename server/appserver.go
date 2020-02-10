package appserver

import (
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const (
	appserverkey  = "/etc/lb/certs/key.pem"
	appservercert = "/etc/lb/certs/cert.pem"
)

var (
	//ignoring developer certificates we are using
	//allows us working with HTTPS without any problems with certificates
	transport = http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	httpclient *http.Client
)

func init() {
	httpclient = &http.Client{Transport: &transport}
}

func populateTemplates() *template.Template {
	result := template.New("templates")
	const basePath = "server/templates"
	template.Must(result.ParseGlob(basePath + "/*.html"))
	return result
}

func Run(lbDiscoUrl string) {

	var loadbalancerURL = flag.String("loadbalancer", lbDiscoUrl, "Address of the load balancer")

	fmt.Printf("listening on https://localhost:3000/. Execute  template home.html by https://localhost:3000/home \n")
	templates := populateTemplates()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//remove slash
		requestedFile := r.URL.Path[1:]
		if requestedFile == "" {
			requestedFile = "index"
		}
		t := templates.Lookup(requestedFile + ".html")
		if t != nil {
			err := t.Execute(w, nil)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//these two are handled automatically by fileserver
	http.Handle("/img/", http.FileServer(http.Dir("server/public")))
	http.Handle("/css/", http.FileServer(http.Dir("server/public")))

	go func() {

		time.Sleep(5 * time.Second)
		fmt.Println("registering; loadbalancer url: " + *loadbalancerURL + "/register?port=3000")

		_, err := httpclient.Get(*loadbalancerURL + "/register?port=3000")
		if err != nil {
			fmt.Println(err)
		}

	}()

	http.ListenAndServeTLS(":3000", appservercert, appserverkey, nil)

	http.Get(*loadbalancerURL + "/unregister?port=3000")
}
