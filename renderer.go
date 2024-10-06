package bananas

import (
	"bufio"
	"bytes"
	"embed"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplData map[string]any

func MustRenderTmpl(
	templateFs embed.FS, templateFile string, data TemplData,
) (outPath string, content *bytes.Buffer, ok bool) {

	ok = strings.HasSuffix(templateFile, ".tmpl")
	if !ok {
		log.Println("skipping file", templateFile, ". file name should end with .tmpl")
		return
	}

	tmplFile, err := templateFs.Open(templateFile)
	if err != nil {
		log.Println("failed to open file", templateFile, "reason:", err)

		return "", nil, false
	}
	defer tmplFile.Close()

	scanner := bufio.NewScanner(tmplFile)

	// Read the first line
	if scanner.Scan() {
		firstLine := scanner.Text()
		if strings.HasPrefix(firstLine, "// out_path:") {
			outPath = strings.TrimSpace(strings.TrimPrefix(firstLine, "// out_path:"))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("failed to scan file", templateFile, "reason:", err)
		return "", nil, false
	}

	// Read the rest of the template
	var tmplContent bytes.Buffer
	for scanner.Scan() {
		tmplContent.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		log.Println("failed to read template", templateFile, "reason:", err)

		return "", nil, false
	}

	// Parse and execute the template
	t, err := template.New(filepath.Base(templateFile)).Parse(tmplContent.String())
	if err != nil {
		log.Println("failed to create template for", templateFile, "reason:", err)
		return "", nil, false
	}

	content = &bytes.Buffer{}
	if err := t.Execute(content, data); err != nil {
		log.Println("failed to render template for", templateFile, "reason:", err)
		return "", nil, false
	}

	return outPath, content, true
}

func WriteFile(filePath string, content *bytes.Buffer) error {
	// Ensure the target directory exists
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); err != nil {
		log.Fatalf("error in setup, %s dir not created. reason: %v\n", dir, err)
	}

	// Write the file
	err := os.WriteFile(filePath, content.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Eeror writing file %s. reason: %v\n", filePath, err)
	}

	return err
}
