package seed

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func writeCV(uploadDir, folder, name, headline string) (string, error) {
	path := folder + "/" + rand.Text() + ".pdf"
	fullPath := filepath.Join(uploadDir, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o750); err != nil {
		return "", err
	}
	return path, os.WriteFile(fullPath, samplePDF(name, headline), 0o640)
}

func samplePDF(name, headline string) []byte {
	escape := strings.NewReplacer(`\`, `\\`, "(", `\(`, ")", `\)`).Replace
	content := fmt.Sprintf("BT /F1 22 Tf 72 760 Td (%s) Tj /F1 13 Tf 0 -28 Td (%s) Tj 0 -40 Td (CV contoh untuk data demo Karsa.) Tj ET", escape(name), escape(headline))
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>",
		fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(content), content),
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
	}
	var document bytes.Buffer
	document.WriteString("%PDF-1.4\n")
	offsets := make([]int, len(objects))
	for index, object := range objects {
		offsets[index] = document.Len()
		fmt.Fprintf(&document, "%d 0 obj\n%s\nendobj\n", index+1, object)
	}
	tableOffset := document.Len()
	fmt.Fprintf(&document, "xref\n0 %d\n0000000000 65535 f \n", len(objects)+1)
	for _, offset := range offsets {
		fmt.Fprintf(&document, "%010d 00000 n \n", offset)
	}
	fmt.Fprintf(&document, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objects)+1, tableOffset)
	return document.Bytes()
}
