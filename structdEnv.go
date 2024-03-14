package structdEnv

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func main() {
}

func Get[T any]() T {
	rm, err := coerceType[T](envMap())
	if err != nil {
		println("error trying to coerce type")
	}
	return rm
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
