package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

/*		Public function: Serialize
 *		---------------------------
 *      Serialize the interface to JSON array of bytes
 *
 *		returns: JSON data in array of bytes
 */
func Serialize(v interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(v)
	return jsonData, err
}

/* 		Public function: CreateRequest
 * 		------------------------------
 *		Create HTTP request
 *
 *		method: HTTP method type (GET, POST, PUT etc)
 *		url: URL endpoint
 *		body: HTTP body
 *
 *		returns: filled http.Request struct.
 */
func CreateRequest(method string, url string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil
	}
	if request != nil {
		request.Header.Add("Content-Type", "application/json; charset=utf-8")
	}
	return request
}

func SendRequest(client *http.Client, req *http.Request) (*http.Response, string, error) {
	// HTTP client send request and get response
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return resp, string(body), nil
}
