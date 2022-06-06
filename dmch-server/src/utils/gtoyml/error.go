package gtoyml

type convertError struct{}

func (e *convertError) Error() string {
	return "gtoyml: convert error"
}
