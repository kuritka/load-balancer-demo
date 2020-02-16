package cache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"lb/common/log"
)

const (
	appserverkey  = "/etc/lb/certs/key.pem"
	appservercert = "/etc/lb/certs/cert.pem"
	cacheKey      = "key"
)

var (
	logger = log.Log
	cache  = make(map[string]*cacheEntry)
	//we don't care how many entries are reading from cache, that's why read Write mutex
	//this mutex can have multiple readers but one writer
	//we expect three servers so no massive traffic, we can use mutexm
	mutex        = sync.RWMutex{}
	maxAgeRegexp = regexp.MustCompile(`maxage=(\d+)`)
	//makes sense  60 sec in PROD because we have to lock cache to inspect all entries
	tickCh = time.Tick(5 * time.Second)
)

type (
	cacheEntry struct {
		data       []byte
		expiration time.Time
	}
)

func Run(port string) {
	fmt.Println("cache started")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getFromCache(w, r)
		case http.MethodPost:
			saveToCache(w, r)
		}
	})

	http.HandleFunc("/invalidate", func(w http.ResponseWriter, r *http.Request) {
		invalidateEntry(w, r)
	})

	go purgeCache()

	http.ListenAndServeTLS(port, appservercert, appserverkey, nil)

	logger.Info().Msgf("server running on port %s", port)
}

func purgeCache() {
	for range tickCh {
		mutex.Lock()
		now := time.Now()

		fmt.Printf("purging cache \n")

		for k, v := range cache {
			if now.Before(v.expiration) {
				fmt.Printf("purging entry %s \n", k)
				delete(cache, k)
			}
		}
		mutex.Unlock()
	}
}

func invalidateEntry(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	key := r.URL.Query().Get(cacheKey)
	fmt.Printf("")
	delete(cache, key)

}

func saveToCache(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	key := r.URL.Query().Get(cacheKey)
	cacheHeader := r.Header.Get("cache-control")

	fmt.Printf("Saving cache entry with key '%s' for %s seconds", key, cacheHeader)

	matches := maxAgeRegexp.FindStringSubmatch(cacheHeader)

	if len(matches) == 2 {
		dur, _ := strconv.Atoi(matches[1])
		data, _ := ioutil.ReadAll(r.Body)
		cache[key] = &cacheEntry{data: data, expiration: time.Now().Add(time.Duration(dur) * time.Second)}

	}
}

func getFromCache(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()
	key := r.URL.Query().Get(cacheKey)

	fmt.Printf("Searching cahce for %s", key)

	if entry, ok := cache[key]; ok {
		fmt.Println("found")
		w.Write(entry.data)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Println("not found")
}
