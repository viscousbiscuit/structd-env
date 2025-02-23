package structdEnv_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	structdEnv "github.com/viscousbiscuit/structd-env"
)

type MyTestType struct {
	FirstName    string
	LastName     string `env:"LAST_NAME"`
	BusinessName string
	Age          uint `env:"AGE"`
	NetWorth     float64
	Active       bool `env:"IS_ACTIVE"`
}

func TestMain(m *testing.M) {

	mockJson := make(map[string]interface{})
	mockJson["FIRST_NAME"] = "Bob"

	fs, _ := json.Marshal(mockJson)

	err := os.WriteFile(".env.json", fs, 0644)
	if err != nil {
		fmt.Println("Unable to create .env.json file")
		os.Exit(1)
	}

	code := m.Run()
	os.Remove(".env.json")

	os.Exit(code)
}

func TestSet(test *testing.T) {

	se, _ := structdEnv.GetInstance[MyTestType]()
	os.Setenv("FIRST_NAME", "Bob")
	se.Flush()

	se.Set(
		MyTestType{
			FirstName:    "Alice",
			LastName:     "",
			BusinessName: "",
			Age:          0,
			NetWorth:     0,
			Active:       false,
		})

	if se.Get().FirstName != "Alice" {
		test.Fail()
	}
}

func TestLower(test *testing.T) {

	se, _ := structdEnv.GetInstance[MyTestType]()
	os.Unsetenv("firstname")
	os.Unsetenv("FirstName")
	os.Setenv("FIRST_NAME", "Bob")
	se.Flush()

	structdEnv := se.Get()

	if structdEnv.FirstName == "" {
		test.Fail()
	}
}

func TestPriority(test *testing.T) {

	se, _ := structdEnv.GetInstance[MyTestType]()
	os.Unsetenv("LAST_NAME")
	os.Unsetenv("lastname")
	os.Setenv("LAST_NAME", "Messerschmitt")
	os.Setenv("lastname", "Kindelberger")

	se.Flush()
	structdEnv := se.Get()
	fmt.Print("Last Name: ", structdEnv.LastName)

	if structdEnv.LastName != "Messerschmitt" {
		test.Fail()
	}
}

func TestFloat(test *testing.T) {

	se, _ := structdEnv.GetInstance[MyTestType]()
	os.Unsetenv("NET_WORTH")
	os.Setenv("NET_WORTH", "123.45")
	se.Flush()

	structdEnv := se.Get()

	if structdEnv.NetWorth != 123.45 {
		test.Fail()
	}
}

func TestUnsigned(test *testing.T) {

	se, _ := structdEnv.GetInstance[MyTestType]()
	os.Unsetenv("AGE")
	os.Setenv("AGE", "-1")
	se.Flush()

	structdEnv := se.Get()

	if structdEnv.Age != 0 {
		test.Fail()
	}
}

func TestBool(test *testing.T) {

	se, _ := structdEnv.GetInstance[MyTestType]()
	os.Unsetenv("IS_ACTIVE")
	os.Setenv("IS_ACTIVE", "true")

	se.Flush()
	structdEnv := se.Get()

	if structdEnv.Active != true {
		test.Fail()
	}
}

func TestBoolNum(test *testing.T) {

	se, _ := structdEnv.GetInstance[MyTestType]()
	os.Unsetenv("IS_ACTIVE")
	os.Setenv("IS_ACTIVE", "1")
	se.Flush()

	structdEnv := se.Get()

	if structdEnv.Active != true {
		test.Fail()
	}
}
