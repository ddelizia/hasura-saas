package hstype

type (
	String  = *string
	Int64   = *int64
	Bool    = *bool
	Float64 = *float64
)

func NewString(s string) String {
	return &s
}

func NewInt64(s int64) Int64 {
	return &s
}

func NewBool(b bool) Bool {
	return &b
}

func NewFloat64(f float64) Float64 {
	return &f
}
