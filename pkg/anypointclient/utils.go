package anypointclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func checkResponse(resp *http.Response, err error) error {
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if resp.StatusCode == 401 {
		return fmt.Errorf("unauthorized access. Please verify that authToken is valid")
	} else if resp.StatusCode > 299 {
		return fmt.Errorf("call failed with status code %d", resp.StatusCode)
	}
	return nil
}

func readResponseAs(resp *http.Response, v interface{}) error {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, &v)
	if err != nil {
		return err
	}
	return nil
}
