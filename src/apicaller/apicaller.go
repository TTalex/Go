package apicaller

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func Callapi(url string) (map[string]interface{}, error){
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var decoded interface{}
	errj := json.Unmarshal(body, &decoded)
	if errj != nil {
		return nil, err
	}
	return decoded.(map[string]interface{}), nil
	
}

func Callapisem(url string, c chan bool) (map[string]interface{}, error){
	c <- true
	defer func() {<-c}()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var decoded interface{}
	errj := json.Unmarshal(body, &decoded)
	if errj != nil {
		return nil, err
	}
	return decoded.(map[string]interface{}), nil
	
}
