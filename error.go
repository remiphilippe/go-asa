package goasa

import (
	"encoding/json"
	"fmt"
)

// ASAMessage  Error message returned by API
// type ASAMessage struct {
// 	Level   string
// 	Details string
// 	Code    string
// }

// ASAError Error returned by API
type ASAError struct {
	Level   string
	Details string
	Code    string
	// Severity string
	// Messages []ASAMessage
}

func (ae ASAError) Error() string {
	return fmt.Sprintf("%s: with code %s and message %s", ae.Level, ae.Code, ae.Details)
}

func parseResponse(bodyText []byte) (err error) {
	//spew.Dump(string(bodyText))
	if len(bodyText) > 0 {
		var v struct {
			Error []ASAError `json:"messages"`
		}

		//log.Print("Response: " + string(bodyText))

		err = json.Unmarshal(bodyText, &v)
		if err != nil {
			return err
		}

		if len(v.Error) == 1 {
			return v.Error[0]
		}

		return fmt.Errorf("we finally found an error with more than 1 message! %+v", v)
	}

	return nil
}
