package receivers

import "go-playground/receivers/auth"

func main() {
	a := auth.NewAPIAuth([]string{"key1", "key2"})
	a.Set([]string{"key3", "key4"})
	a.Authenticate("key3")
}
