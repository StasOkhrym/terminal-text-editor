package main

import (
	"fmt"
	"os"

	"text_editor/editor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Give me <file_path>")
		return
	}

	filePath := os.Args[1]

	e, err := editor.NewEditor(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer e.Close()

	if err := e.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
