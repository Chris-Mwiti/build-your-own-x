package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main(){
	//map marshaling to json
	data := map[string]interface{} {
		"intValue": 1234,
		"boolValue": true,
		"stringValue": "hello!",
		"dateValue": time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.UTC),
		"objectValue": map[string]interface{}{
			"arrayValue": []int{1, 2,3,4},
		},
		"nilStringValue": nil,
		"nilIntValue": nil,

	}
	jsonData, err := json.Marshal(data)

	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	var myInt int = 5
	var ptrInt *int = &myInt

	//struct marshaling 
	type myObject struct {
		ArrayValue []int `json:"arrayValue"`
	}
	type myStruct struct {
		IntValue int `json:"intValue"` 
		BoolValue bool `json:"booValue"`
		StringValue string `json:"stringValue"`
		DateValue time.Time `json:"dateValue"`
		ObjectValue *myObject `json:"objectValue"`
		NullIntValue *int `json:"nullIntValue,omitempty"`
		NullStringValue *string `json:"nullStringValue,omitempty"`
	}
	 
	structData := &myStruct{
		IntValue: 1235,
		BoolValue: true,
		StringValue: "hello world!",
		DateValue: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.UTC),
		ObjectValue: &myObject{
			ArrayValue: []int{1,2,3,4},
		},
		NullStringValue: nil,
		NullIntValue: ptrInt,
	}
	structJsonData, err := json.Marshal(structData)
	fmt.Printf("json data: %s\n", jsonData)
	fmt.Printf("struct json data: %s\n", structJsonData)
}