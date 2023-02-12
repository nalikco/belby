package request

import (
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	Method   string
	URL      string
	Body     io.Reader
	Query    map[string]string
	Headers  map[string]string
	Callback func(body io.ReadCloser) ([]byte, error)
}

func SendRequest(r Request) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		return []byte{}, err
	}

	for key, param := range r.Headers {
		req.Header.Set(key, param)
	}

	if len(r.Query) > 0 {
		q := req.URL.Query()
		for key, param := range r.Query {
			q.Add(key, param)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	return r.Callback(resp.Body)
}
