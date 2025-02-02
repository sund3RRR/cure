package types

type Sequence string

const (
	UnderlinedSeq Sequence = "\033[4m"
	ResetSeq      Sequence = "\033[0m"
)
