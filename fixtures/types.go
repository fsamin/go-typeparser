package main

func main() {}

//TypeA Comment
//metadata
type TypeA struct {
	FieldA int  `tag:"value,option" json:"-"`
	FieldB bool `tag:"value,option"`
}

//InterfaceA Comment
type InterfaceA interface {
	foo(string) (string, error)
	bar(...int) (string, error)
}
