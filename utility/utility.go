package utility

import (
	"strings"
)

func Quote(str0 string) string {
	str1 := strings.Replace(str0, "\\", `\\`, -1)
	str2 := strings.Replace(str1, "\n", `\n`, -1)
	str3 := strings.Replace(str2, "\t", `\t`, -1)
	str4 := strings.Replace(str3, "\"", `\"`, -1)
	return `"`+str4+`"`
}
