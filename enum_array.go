package main

var _EnumArray = []Enum{
	Enum1,
	Enum2,
	Enum3,
	Enum5,
}

func (Enum) EnumArray() []Enum {
	return _EnumArray
}
