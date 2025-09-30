package executor

import "redis-repo/internal/data_structure"

var dictStore *data_structure.Dict

func init() {
	dictStore = data_structure.NewDict()
}
