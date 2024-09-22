package files

import (
	"bufio"
	"fmt"
	"os"
)

func WriteToFile(filename string, data []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range data {
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}
