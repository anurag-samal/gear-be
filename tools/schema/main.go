package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const migrationsDir = "system/migrations"

func main() {
	output := flag.String("o", "system/schema/schema.sql", "output file path (use '-' for stdout)")
	flag.Parse()

	schema := generate()

	if *output == "-" {
		os.Stdout.WriteString(schema)
		return
	}

	if err := os.MkdirAll(filepath.Dir(*output), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(*output, []byte(schema), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("Schema written to %s\n", *output)
}

func generate() string {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		panic(err)
	}

	slices.SortFunc(entries, func(a, b os.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	var out strings.Builder

	out.WriteString("-- ============================================================\n")
	out.WriteString("-- Gear Database Schema\n")
	out.WriteString("-- Auto-generated from system/migrations/ — do not edit directly\n")
	out.WriteString("-- ============================================================\n\n")

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
		if err != nil {
			panic(err)
		}

		content := strings.TrimSpace(string(data))
		if content == "" {
			continue
		}

		label := labelFromFilename(entry.Name())

		out.WriteString("-- ============================================================\n")
		out.WriteString(fmt.Sprintf("-- %s\n", label))
		out.WriteString("-- ============================================================\n\n")
		out.WriteString(content)
		out.WriteString("\n\n")
	}

	return out.String()
}

func labelFromFilename(name string) string {
	name = strings.TrimSuffix(name, ".sql")
	parts := strings.SplitN(name, "_", 2)
	label := name
	if len(parts) == 2 {
		label = parts[1]
	}
	label = strings.ReplaceAll(label, "_", " ")
	label = strings.ToUpper(label[:1]) + label[1:]
	return label
}
