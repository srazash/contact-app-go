package main

import (
	"contactapp/internal/contact"
	"fmt"
)

func main() {
	ptrDB := contact.All()
	fmt.Printf("DB len: %d\n", len(*ptrDB))
}
