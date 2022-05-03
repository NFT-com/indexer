package filters

import (
	"fmt"
	"strings"
)

func Eq(field string, value string) string {
	return fmt.Sprintf("%s = %s", field, value)
}

func In(field string, values ...string) string {
	value := strings.Join(values, ",")
	value = "(" + value + ")"
	return fmt.Sprintf("%s IN %s", field, value)
}
