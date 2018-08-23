package goasa

import (
	"encoding/json"
	"fmt"
)

// ASAMessage  Error message returned by API
type ASAMessage struct {
	Level   string
	Details string
	Code    string
}

// ASAError Error returned by API
type ASAError struct {
	Severity string
	Messages []ASAMessage
}

func (ae ASAError) Error() string {
	return fmt.Sprintf("%s: with messages %+v", ae.Severity, ae.Messages)
}

func parseResponse(bodyText []byte) (err error) {
	//spew.Dump(string(bodyText))
	if len(bodyText) > 0 {
		var v struct {
			Error []ASAMessage `json:"messages"`
		}

		//log.Print("Response: " + string(bodyText))

		err = json.Unmarshal(bodyText, &v)
		if err != nil {
			return err
		}

		aErr := new(ASAError)
		aErr.Messages = v.Error
		//TODO: need to populate the asaError with the highest severity found

		return aErr
	}

	return nil
}
