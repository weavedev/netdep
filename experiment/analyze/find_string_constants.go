package analyze

import (
	"fmt"
	"go/constant"
	"golang.org/x/tools/go/ssa"
)

func findReferenceValues(referrers *[]ssa.Instruction, visited []*ssa.Value, depth int) {
	if referrers == nil || len(*referrers) == 0 || depth > 16 {
		return
	}

	for _, ref := range *referrers {
		switch instr := ref.(type) {
		case *ssa.Call:
			for _, arg := range instr.Call.Args {
				findStringConstants(arg, visited)
			}
		case *ssa.Return:
			findReferenceValues(instr.Referrers(), visited, depth+1)
			break
		case ssa.Value:
			findStringConstants(instr, visited)
			break
		default:
			break
		}
	}
}

func findStringConstants(value ssa.Value, visited []*ssa.Value) {
	//newFound := make([]string, 0)

	if len(visited) > 32 {
		return
	}

	for _, v := range visited {
		if v == &value {
			return
		}
	}

	switch val := value.(type) {
	case *ssa.Const:
		switch val.Value.Kind() {
		case constant.String:
			constVal := constant.StringVal(val.Value)
			//append(newFound, constant.StringVal(val.Value))
			fmt.Println("Found const: " + constVal)
			return
		}
	case *ssa.Extract:
		findStringConstants(val.Tuple, append(visited, &value))
	default:
		break
	}

	findReferenceValues(value.Referrers(), append(visited, &value), len(visited))
}
