package main

import (
	"contactapp/internal/contact"
	"fmt"
)

func main() {
	ptrDB := contact.All()
	fmt.Printf("DB len: %d\n", len(*ptrDB))

	//contact.CreateContact("Ryan", "Shaw-Harrison", "ryan@mail.local", "+44 (0) 1234 567890")
	contact.LoadDB()
	fmt.Printf("DB loaded, len: %d\n", len(*ptrDB))
	fmt.Printf("DB: %v\n", *ptrDB)

	for _, c := range *ptrDB {
		fmt.Printf("\tID: %d\n", c.Id)
		fmt.Printf("\tName: %s, %s\n", c.Last, c.First)
		fmt.Printf("\tEmail: %s\n", c.Email)
		fmt.Printf("\tPhone: %s\n", c.Phone)
	}
}
