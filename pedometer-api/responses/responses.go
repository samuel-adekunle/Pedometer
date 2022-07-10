package responses

import "encoding/json"

func JsonResponseBody(message string, data interface{}) (string, error) {
	jsonData, err := json.Marshal(map[string]interface{}{
		"message": message,
		"data":    data,
	})
	return string(jsonData), err
}

func JsonResponseHeader() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin": "*",
	}
}
