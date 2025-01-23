package main

import (
 "encoding/json"
 "fmt"
)

func parsing() {
	var jsonData = `
		{
			"intValue": 1234,
			"boolValue": true,
			"stringValue": "hello!",
			"dateValue": "2025-01-02T09:10:00Z",
			"objectValue": {
				"arrayValue": [1,2,3,4]	
			},
			"nullStringValue": null,
			"nullIntValue": null
		}	
	`

	var data *myStruct 
	err := json.Unmarshal([]byte(jsonData),&data)
	if err != nil {
		fmt.Printf("could not unmarshal json: %s\n", err)
		return
	}
	fmt.Printf("json map: %#v\n", data)
	fmt.Printf("datevalue: %#v\n", data.DateValue)
	fmt.Printf("objectvalue: %#v\n", data.ObjectValue)


}