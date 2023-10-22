package problems

type CannotPerformOperationError struct {
	Reason string
}

func (c CannotPerformOperationError) Error() string {
	return c.Reason
}

type BadArgumentsError struct {
	Reason string
}

func (c BadArgumentsError) Error() string {
	return c.Reason
}
