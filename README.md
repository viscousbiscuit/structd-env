# structd-env

A simple package to manage environment variables in a structured way.

## Usage

```
type MyType struct {
	FirstName    string
	LastName     string `env:"LAST_NAME"`
	BusinessName string
	Age          uint `env:"AGE"`
	NetWorth     float64
	Active       bool `env:"IS_ACTIVE"`
}

func Main() {
    structdEnv := Make[MyTestType]()
    println(structdEnv.FirstName)
}
```

By default, a best effort attempt to match the environment variable to the struct field will be made. If two environment variables exist with the same name, but different casing use the `env` tag 
to specify the environment variable name to use for that field.

Adding `env` tags to your struct will allow you to specify the environment variable name to use for that field. If no `env` tag is provided, the field name will be used as the environment variable name.

