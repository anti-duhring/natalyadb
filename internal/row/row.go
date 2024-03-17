package row

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/anti-duhring/natalyadb/pkg/utils"
)

func SerializeRow(s interface{}) string {
	t := reflect.TypeOf(s)
	index := ""
	serialized := ""

	cursor := 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := strings.ToLower(field.Name)
		value := fmt.Sprint(reflect.ValueOf(s).Field(i).Interface())

		index += utils.BytesToString(utils.IntToBytes(cursor)) + utils.BytesToString(utils.IntToBytes(len(value)))

		cursor += len(name)
		serialized += fmt.Sprint(value)
	}

	return index + serialized
}
