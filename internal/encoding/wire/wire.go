package wire

type Number int32

const (
	MinValidNumber      Number = 1
	FirstReservedNumber Number = 19000
	LastReservedNumber  Number = 19999
	MaxValidNumber      Number = 1<<29 - 1
)

func (n Number) IsValid() bool {
	return MinValidNumber <= n && n < FirstReservedNumber ||
		LastReservedNumber < n && n <= MaxValidNumber
}

type Type int8

const ()
