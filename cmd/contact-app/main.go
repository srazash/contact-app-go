package main

import (
	"contactapp/internal/contact"
	"fmt"
)

func main() {
	fmt.Printf("DB len: %d\n", len(contact.DB))
}
