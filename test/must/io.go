package must

import "io"

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		panic(err)
	}
}
