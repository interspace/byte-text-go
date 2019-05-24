package extract

import (
	"fmt"
	"io/ioutil"
	"testing"

	goyaml "gopkg.in/yaml.v1"
)

func TestMentions(t *testing.T) {
	contents, err := ioutil.ReadFile(extractYmlPath)
	if err != nil {
		t.Errorf("Error reading extract.yml: %v", err)
		t.FailNow()
	}

	var conformance = &Conformance{}
	err = goyaml.Unmarshal(contents, &conformance)
	if err != nil {
		t.Errorf("Error parsing extract.yml: %v", err)
		t.FailNow()
	}

	mentionTests, ok := conformance.Tests["mentions"]
	if !ok {
		t.Errorf("Conformance file did not contain 'mentions' key")
		t.FailNow()
	}

	for _, test := range mentionTests {
		result := MentionedScreenNames(test.Text)

		expected, ok := test.Expected.([]interface{})
		if !ok {
			fmt.Printf("e: %#v\n", test)
			t.Errorf(
				"Expected value in conformance file was not a list. Test name: %s.\n",
				test.Description,
			)
			t.FailNow()
		}

		if len(result) != len(expected) {
			t.Errorf(
				"Wrong number of entities returned for text [%s]. Expected:%v Got:%v.\n",
				test.Text,
				expected,
				result,
			)
			continue
		}

		for n, e := range expected {
			actual := result[n]
			if actual.screenName != e {
				t.Errorf(
					"ExtractMentionedScreenNames returned incorrect value for test: [%s]. Expected:[%s] Got:[%s]\n",
					test.Text,
					e,
					actual.Text,
				)
			}

			if actual.Type != Mention {
				t.Errorf(
					"ExtractMentionedScreenNames returned entity with wrong type. Expected:Mention Got:%v",
					actual.Type,
				)
			}
		}
	}
}

func TestMentionsWithIndices(t *testing.T) {
	contents, err := ioutil.ReadFile(extractYmlPath)
	if err != nil {
		t.Errorf("Error reading extract.yml: %v", err)
		t.FailNow()
	}

	var conformance = &Conformance{}
	err = goyaml.Unmarshal(contents, &conformance)
	if err != nil {
		t.Errorf("Error parsing extract.yml: %v", err)
		t.FailNow()
	}

	mentionTests, ok := conformance.Tests["mentions_with_indices"]
	if !ok {
		t.Errorf("Conformance file did not contain 'mentions_with_indices' key")
		t.FailNow()
	}

	for _, test := range mentionTests {
		result := MentionedScreenNames(test.Text)

		expected, ok := test.Expected.([]interface{})
		if !ok {
			fmt.Printf("e: %#v\n", test)
			t.Errorf(
				"Expected value in conformance file was not a list. Test name: %s.\n",
				test.Description,
			)
			t.FailNow()
		}

		if len(result) != len(expected) {
			t.Errorf(
				"Wrong number of entities returned for text [%s]. Expected:%v Got:%v.\n",
				test.Text,
				expected,
				result,
			)
			continue
		}

		for n, e := range expected {
			actual := result[n]
			expectedMap, ok := e.(map[interface{}]interface{})
			if !ok {
				t.Errorf(
					"Expected value was not a map. Test name: %s\n",
					test.Description,
				)
				continue
			}

			mention, ok := expectedMap["screen_name"]
			if !ok {
				t.Errorf(
					"Expected value did not contain screen_name. Test name: %s\n",
					test.Description,
				)
				continue
			}

			if actual.screenName != mention {
				t.Errorf(
					"ExtractMentionedScreenNames returned incorrect value for test: [%s]. Expected:[%s] Got:[%s]\n",
					test.Text,
					mention,
					actual.screenName,
				)
			}

			indices, ok := expectedMap["indices"]
			if !ok {
				t.Errorf(
					"Expected value did not contain indices. Test name: %s\n",
					test.Description,
				)
				continue
			}

			indicesList := indices.([]interface{})
			if len(indicesList) != 2 {
				t.Errorf(
					"Indices did not contain 2 values. Test name: %s\n",
					test.Description,
				)
				continue
			}

			if indicesList[0] != actual.Range.Start ||
				indicesList[1] != actual.Range.Stop {
				t.Errorf(
					"ExtractMentionedScreenNames did not return correct indices [%s]. Expected:(%d, %d) Got:%s)",
					test.Text,
					indicesList[0],
					indicesList[1],
					actual.Range,
				)
			}
		}
	}
}
