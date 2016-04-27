package main

import "github.com/bndr/gotabulate"

// BuildTable builds a table to output
func BuildTable(headers []string, data [][]string) string {

	t := gotabulate.Create(data)

	// Set the Headers (optional)
	t.SetHeaders(headers)

	// Set the Empty String (optional)
	t.SetEmptyString("")

	// Set Align (Optional)
	t.SetAlign("left")

	// Print the result: grid, or simple
	return t.Render("grid")
}
