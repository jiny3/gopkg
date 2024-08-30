package structx

import "github.com/fatih/structs"

func Struct2Map(in any) map[string]interface{} {
	return structs.Map(in)
}
