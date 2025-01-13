package httptest

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[missingKey-1]
	_ = x[missingValue-2]
	_ = x[unkownParameter-3]
	_ = x[unknownValue-4]
	_ = x[valueMismatch-5]
}

const _diffType_name = "PARAMETER MISSINGMISSING VALUEUNEXPECTED PARAMETERUNKNOWN VALUEVALUE MISMATCH"

var _diffType_index = [...]uint8{0, 17, 30, 50, 63, 77}

func (i diffType) String() string {
	i -= 1
	return _diffType_name[_diffType_index[i]:_diffType_index[i+1]]
}
