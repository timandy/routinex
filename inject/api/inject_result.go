package api

type InjectResult struct {
	ReplaceFiles map[int]string
	ExtraFiles   []string
}

func NewInjectResult() *InjectResult {
	return &InjectResult{ReplaceFiles: map[int]string{}, ExtraFiles: []string{}}
}
