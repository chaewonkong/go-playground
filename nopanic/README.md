# nopanic
A struct has a field of an interface type. At runtime, the field is assigned nil. However, calling a method on this nil interface field does not result in a panic!