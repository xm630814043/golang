package marshal

import (
	"encoding/json"
	"fmt"
)

type TestStruct struct {
	Type int
	Body json.RawMessage
}

type Person struct {
	Name string
	Age  int
}

type Worker struct {
	Name string
	Job  string
}

func RawMessage() {
	input := `
       {
        "Type": 1,
        "Body":{ 
            "Name":"ff",
            "Age" : 19
         }
    }`
	ts := TestStruct{}
	if err := json.Unmarshal([]byte(input), &ts); err != nil {
		panic(err)
	}
	switch ts.Type {
	case 1:
		var p Person
		if err := json.Unmarshal(ts.Body, &p); err != nil {
			panic(err)
		}
		fmt.Println(p)
	case 2:
		var w Worker
		if err := json.Unmarshal(ts.Body, &w); err != nil {
			panic(err)
		}
		fmt.Println(w)
	}
}
