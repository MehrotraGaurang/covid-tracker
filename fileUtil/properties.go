package fileutil

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var AppConfigProperties map[string]string = make(map[string]string)

func ReadPropertiesFile(filename string) error {

	if len(filename) == 0 {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				AppConfigProperties[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
