package validate

import (
	"io/ioutil"
	"testing"

	goyaml "gopkg.in/yaml.v1"
)

func TestTextIsValid(t *testing.T) {
	contents, err := ioutil.ReadFile(validateYmlPath)
	if err != nil {
		t.Errorf("Error reading validate.yml: %v", err)
		t.FailNow()
	}

	var testData map[interface{}]interface{}
	err = goyaml.Unmarshal(contents, &testData)
	if err != nil {
		t.Fatalf("error unmarshaling data: %v\n", err)
	}

	tests, ok := testData["tests"]
	if !ok {
		t.Errorf("Conformance file was not in expected format.")
		t.FailNow()
	}

	textTests, ok := tests.(map[interface{}]interface{})["texts"]
	if !ok {
		t.Errorf("Conformance file did not contain text tests")
		t.FailNow()
	}

	for _, testCase := range textTests.([]interface{}) {
		test := testCase.(map[interface{}]interface{})
		text, _ := test["text"]
		description, _ := test["description"]
		expected, _ := test["expected"]

		actual := TextIsValid(text.(string), ValidationArgs{canBeEmpty: false, maxLength: 140})
		if actual != expected {
			t.Errorf(
				"TextIsValid returned incorrect value for test [%s]. Expected:%v Got:%v",
				description,
				expected,
				actual,
			)
		}
	}
}
