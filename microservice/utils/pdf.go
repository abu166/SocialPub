package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePDF generates a PDF receipt directly using gofpdf.
func GeneratePDF(transactionID string, details map[string]string) (string, error) {
	// Ensure the receipts directory exists
	err := os.MkdirAll("receipts", os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create receipts directory: %v", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Add content to the PDF
	pdf.Cell(40, 10, "Company Name: My Awesome Project")
	pdf.Ln(10)
	pdf.Cell(40, 10, "Transaction ID: "+transactionID)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Order Date: "+details["OrderDate"])
	pdf.Ln(10)
	pdf.Cell(40, 10, "Product/Service: "+details["ProductName"])
	pdf.Ln(10)
	pdf.Cell(40, 10, "Price Per Unit: "+details["PricePerUnit"])
	pdf.Ln(10)
	pdf.Cell(40, 10, "Quantity: "+details["Quantity"])
	pdf.Ln(10)
	pdf.Cell(40, 10, "Total Amount: "+details["TotalAmount"])
	pdf.Ln(10)
	pdf.Cell(40, 10, "Customer Name: "+details["CustomerName"])
	pdf.Ln(10)
	pdf.Cell(40, 10, "Payment Method: "+details["PaymentMethod"])

	// Save the PDF
	filePath := filepath.Join("receipts", fmt.Sprintf("%s.pdf", transactionID))
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to save PDF: %v", err)
	}

	return filePath, nil
}
