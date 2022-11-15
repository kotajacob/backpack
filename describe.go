package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// description returns the description of an item.
func (b backpack) description(item string) string {
	descriptions, err := b.loadDescriptions()
	if err != nil {
		log.Printf("error loading descriptions: %v\n", err)
		return FatalMessage
	}
	return displayName(item, 1) + ": " + descriptions[normalizeName(item)]
}

// setDescription updates the description of an item.
func (b backpack) setDescription(item, description string) string {
	descriptions, err := b.loadDescriptions()
	if err != nil {
		log.Printf("error loading descriptions: %v\n", err)
		return FatalMessage
	}

	descriptions[normalizeName(item)] = description
	err = b.storeDescriptions(descriptions)
	if err != nil {
		log.Printf("error storing descriptions: %v\n", err)
		return FatalMessage
	}
	return "Updated description of " + item + "."
}

// loadDescriptions returns a mapping of items to descriptions.
func (b backpack) loadDescriptions() (map[string]string, error) {
	descriptions := make(map[string]string)
	file, err := os.Open(filepath.Join(b.dir, "descriptions.kv"))
	if errors.Is(err, fs.ErrNotExist) {
		return descriptions, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		a := strings.Split(line, "=")
		if len(a) != 2 {
			return nil, fmt.Errorf("invalid description: %v", line)
		}
		descriptions[a[0]] = a[1]
	}

	return descriptions, nil
}

// storeDescriptions stores a mapping of items to descriptions.
func (b backpack) storeDescriptions(descriptions map[string]string) error {
	var buf bytes.Buffer
	for item, description := range descriptions {
		buf.WriteString(item)
		buf.WriteString("=")
		buf.WriteString(description)
		buf.WriteString("\n")
	}
	return os.WriteFile(
		filepath.Join(b.dir, "descriptions.kv"),
		buf.Bytes(),
		0777,
	)
}
