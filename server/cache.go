package appserver

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
)

func getFromCache(key string) (io.ReadCloser, bool) {
	resp, err := http.Get(cashUrl + "/?key=" + key)
	if err != nil {
		return nil, false
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}
	logger.Info().Msgf("getting %s from cache", key)
	return resp.Body, true
}

func saveToCache(key string, duration int64, data []byte) {

	req, err := http.NewRequest(http.MethodPost, cashUrl+"/?key="+key, bytes.NewBuffer(data))

	if err != nil {
		logger.Err(err).Msgf("unable to Save %s into cache ", key)
	}

	req.Header.Add("cache-control", "maxage="+strconv.FormatInt(duration, 10))

	http.DefaultClient.Do(req)

	logger.Info().Msgf("'%s' save to cache", key)
}

func invalidateCache(key string) {
	http.Get(cashUrl + "/?key=" + key)
}
