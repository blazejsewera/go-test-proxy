package must

import (
	"fmt"
)

func Succeed[T any](result T, err error) T {
	if err != nil {
		panic(fmt.Errorf("must succeed: unexpected error: %s", err))
	}
	return result
}
