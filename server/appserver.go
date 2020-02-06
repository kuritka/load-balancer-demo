package appserver

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func populateTemplates() *template.Template {
	result := template.New("templates")
	const basePath = "server/templates"
	template.Must(result.ParseGlob(basePath + "/*.html"))
	return result
}

func Run(lbDiscoUrl string) {

	var loadbalancerURL = flag.String("loadbalancer", lbDiscoUrl, "Address of the load balancer")

	fmt.Printf("listening on http://localhost:8000. Execute  template home.html by http://localhost:8000/home \n")
	templates := populateTemplates()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//remove slash
		requestedFile := r.URL.Path[1:]
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

	go http.ListenAndServe(":8000", nil)

	//we must ensure that   http.Get(*loadbalancerURL+"/rigister?port=3000") will be called after server start but
	//we cannot run it one by one as ListenAndServe is blocking function .
	//would be done by donChannel  but in this case I'll wait 1 sec
	time.Sleep(1 * time.Second)
	http.Get(*loadbalancerURL + "/rigister?port=3000")

	log.Println("server started, press <ENTER> to exit")
	_, _ = fmt.Scanln()

	http.Get(*loadbalancerURL + "/unregister?port=3000")
}