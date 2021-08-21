package utility

import (
	"bytes"
	"strings"
)

func Quote(str0 string) string {
	str1 := strings.Replace(str0, "\\", `\\`, -1)
	str2 := strings.Replace(str1, "\n", `\n`, -1)
	str3 := strings.Replace(str2, "\t", `\t`, -1)
	str4 := strings.Replace(str3, "\"", `\"`, -1)
	return `"`+str4+`"`
}

func FunctionString(args []string, body []string) string {
	var out bytes.Buffer
	out.WriteString("(")
	argsLen := len(args)
	for i, name := range args {
		out.WriteString(name)
		if i+1 != argsLen {
			out.WriteString(", ")
		}
	}
	out.WriteString(") {\n")
	for _, stmt := range body {
		tabStmt := strings.Replace(stmt, "\n", "\n\t", -1)
		out.WriteString("\t")
		out.WriteString(tabStmt)
		out.WriteString("\n")
	}
	out.WriteString("}")
	return out.String()
}