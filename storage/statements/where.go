package statements

import (
	"fmt"
	"strings"
)

func Eq(field string, value string) string {
	return fmt.Sprintf("%s = %s", field, value)
}

func In(field string, values ...string) string {
	value := strings.Join(values, ",")
	return fmt.Sprintf("%s IN (%s)", field, value)
}

func NotEq(field string, value string) string {
	return fmt.Sprintf("%s != %s", field, value)
}

func NotIn(field string, values ...string) string {
	value := strings.Join(values, ",")
	return fmt.Sprintf("%s NOT IN (%s)", field, value)
}
