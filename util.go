package typeparser

var (
	Verbose bool
	Logger  func(format string, args ...interface{})
)

func print(format string, args ...interface{}) {
	if Verbose && Logger != nil {
		Logger(format, args...)
	}
}

func debug(i interface{}) {
	print("%T %+v", i, i)
}

// List is a magic list of string
type List []string

// Index returns the first index of the target string `t`, or
// -1 if no match is found.
func (vs List) Index(t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Has returns `true` if the target string t is in the
// slice.
func (vs List) Has(t string) bool {
	return vs.Index(t) >= 0
}

// Any returns `true` if one of the strings in the slice
// satisfies the predicate `f`.
func (vs List) Any(f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// All returns `true` if all of the strings in the slice
// satisfy the predicate `f`.
func (vs List) All(f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

// Filter returns a new slice containing all strings in the
// slice that satisfy the predicate `f`.
func (vs List) Filter(f func(string) bool) List {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return List(vsf)
}

// Map returns a new slice containing the results of applying
// the function `f` to each string in the original slice.
func (vs List) Map(f func(string) string) List {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return List(vsm)
}

// Explode returns a new slice containing the results of applying
// the function `f` to each string in the original slice.
func (vs List) Explode(f func(string) []string) List {
	vsm := []string{}
	for _, v := range vs {
		vsm = append(vsm, f(v)...)
	}
	return List(vsm)
}
