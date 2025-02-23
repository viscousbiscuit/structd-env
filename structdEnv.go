package structdEnv

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var (
	instance any
	once     sync.Once
)

type env[T any] struct {
	value *T
}

// Creates a new singleton instance of the struct
func GetInstance[T any]() (env[T], error) {
	var er error

	once.Do(func() {
		em := getEnvMap()
		rm, err := coerceType[T](em)
		if err != nil {
			er = err
		} else {
			e := env[T]{}
			e.value = &rm
			instance = e
		}
	})

	if er != nil {
		return env[T]{}, er
	}

	return instance.(env[T]), nil
}

// Returns the current environment variables
func (e *env[T]) Get() T {
	return *e.value
}

// Returns the current environment variables
// and updates the cached struct
func (e *env[T]) Flush() (T, error) {
	em := getEnvMap()
	rm, err := coerceType[T](em)
	if err != nil {
		return *e.value, err
	}

	e.value = &rm
	return *e.value, nil
}

// Replaces the current instance of the struct
// with the one provided
func (e *env[T]) Set(val T) T {
	e.value = &val
	return *e.value
}

func coerceType[T any](envMap map[string]string) (T, error) {

	var t T
	tType := reflect.TypeOf(t)
	elem := reflect.ValueOf(&t).Elem()

	if elem.Kind() != reflect.Struct {
		return t, errors.New("not a struct")
	}

	envKeyMap := make(map[string]string)
	for k := range envMap {
		envKeyMap[normalizeKey(k)] = k
	}

	structFieldMap := make(map[string]reflect.StructField)

	for i := 0; i < tType.NumField(); i++ {
		tt := tType.Field(i)
		envTag := tt.Tag.Get("env")
		keyName := normalizeKey(tt.Name)
		if envTag != "" {
			keyName = envTag
		} else {
			mk, ok := envKeyMap[normalizeKey(tt.Name)]
			if ok {
				keyName = mk
			}
		}
		structFieldMap[keyName] = tt
	}

	for k, v := range envMap {
		msf, ok := structFieldMap[k]
		if ok {
			matchedField := elem.FieldByName(msf.Name)
			if matchedField.CanSet() {
				kind := msf.Type.Kind()
				setValue(&matchedField, &v, kind)
			}
		}
	}
	return t, nil
}

func setValue(val *reflect.Value, envVal *string, kind reflect.Kind) {
	switch kind {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		i, err := strconv.Atoi(*envVal)
		if err != nil {
			break
		}
		val.SetInt(int64(i))
		break
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		i, err := strconv.Atoi(*envVal)
		if err != nil {
			break
		}
		if i < 0 {
			i = 0
		}
		val.SetUint(uint64(i))
		break
	case reflect.Float32:
		f, err := strconv.ParseFloat(*envVal, 32)
		if err != nil {
			break
		}
		val.SetFloat(f)
		break

	case reflect.Float64:
		f, err := strconv.ParseFloat(*envVal, 64)
		if err != nil {
			break
		}
		val.SetFloat(f)
		break
	case reflect.Bool:
		tStr := strings.TrimSpace(*envVal)
		if tStr == "1" || tStr == "true" {
			val.SetBool(true)
		} else if tStr == "0" || tStr == "false" {
			val.SetBool(false)
		}
		break
	case reflect.String:
		val.SetString(*envVal)
		break
	}
}

func getEnvMap() map[string]string {

	em := envMap()
	fm := fileEnvMap()
	if len(fm) > 0 {
		for k, v := range fm {
			sv := fmt.Sprintf("%v", v)
			em[k] = sv
		}
	}
	return em
}

func fileEnvMap() map[string]interface{} {
	if _, err := os.Stat(".env.json"); os.IsNotExist(err) {
		return make(map[string]interface{})
	}

	envKeyMap := make(map[string]interface{})
	envFile, err := os.ReadFile(".env.json")

	if err != nil {
		return envKeyMap
	}

	err = json.Unmarshal(envFile, &envKeyMap)
	if err != nil {
		return envKeyMap
	}

	return envKeyMap
}

func envMap() map[string]string {
	envKeyMap := make(map[string]string)
	for _, e := range os.Environ() {
		envVal := strings.Split(e, "=")
		key := &envVal[0]
		val := &envVal[1]

		envKeyMap[*key] = *val
	}
	return envKeyMap
}

func normalizeKey(key string) string {
	rmSnake := strings.Replace(key, "_", "", -1)
	rmDash := strings.Replace(rmSnake, "-", "", -1)
	rmUpper := strings.ToLower(rmDash)
	return rmUpper
}
