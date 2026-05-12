package shared

import (
	"fmt"
	"log"
)

func PrintJSON(label, value string) {
	if label != "" {
		fmt.Printf("== %s ==\n", label)
	}
	fmt.Println(value)
	fmt.Println()
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
