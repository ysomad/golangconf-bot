package schedule

import (
	"fmt"
)

type CSVError struct {
	cause error
	row   int
	col   int
}

func newCSVError(cause error, row, col int) CSVError {
	return CSVError{
		cause: cause,

		// in csv counting rows and columns from 1
		row: row + 1,
		col: col + 1,
	}
}

func (e CSVError) Error() string {
	return fmt.Sprintf("[ряд %d, столбец %d] %s", e.row, e.col, e.cause.Error())
}
