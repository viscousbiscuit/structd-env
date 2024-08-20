package structdEnv

import (
	"os"
	"testing"
)

type MyTestType struct {
	FirstName    string
	LastName     string `env:"LAST_NAME"`
	BusinessName string
	Age          uint `env:"AGE"`
	NetWorth     float64
	Active       bool `env:"IS_ACTIVE"`
}

func TestLower(test *testing.T) {

	os.Unsetenv("firstname")
	os.Unsetenv("FirstName")
	os.Setenv("FIRST_NAME", "Bob")

	structdEnv := Get[MyTestType]()

	if structdEnv.FirstName == "" {
		test.Fail()
	}
}

func TestPriority(test *testing.T) {

	os.Unsetenv("LAST_NAME")
	os.Unsetenv("lastname")

	os.Setenv("LAST_NAME", "Messerschmitt")
	os.Setenv("lastname", "Kindelberger")

	structdEnv := Get[MyTestType]()

	if structdEnv.LastName != "Messerschmitt" {
		test.Fail()
	}
}

func TestFloat(test *testing.T) {

	os.Unsetenv("NET_WORTH")
	os.Setenv("NET_WORTH", "123.45")

	structdEnv := Get[MyTestType]()

	if structdEnv.NetWorth != 123.45 {
		test.Fail()
	}
}

func TestUnsigned(test *testing.T) {

	os.Unsetenv("AGE")
	os.Setenv("AGE", "-1")

	structdEnv := Get[MyTestType]()

	if structdEnv.Age != 0 {
		test.Fail()
	}
}

func TestBool(test *testing.T) {

	os.Unsetenv("IS_ACTIVE")
	os.Setenv("IS_ACTIVE", "true")

	structdEnv := Get[MyTestType]()

	if structdEnv.Active != true {
		test.Fail()
	}
}

func TestBoolNum(test *testing.T) {

	os.Unsetenv("IS_ACTIVE")
	os.Setenv("IS_ACTIVE", "1")

	structdEnv := Get[MyTestType]()

	if structdEnv.Active != true {
		test.Fail()
	}
}
