package main

import "github.com/bndr/gotabulate"

// PrintInstances prints instances
func PrintInstances(instances []InstanceInfo) string {
	var data [][]string

	for _, instance := range instances {
		s := []string{
			instance.Name,
			instance.PrivateIP,
			instance.PublicIP,
			instance.InstanceType,
			instance.InstanceState,
		}
		data = append(data, s)
	}

	headers := []string{"Name", "Private IP", "Public IP", "Type", "State"}

	return buildTable(headers, data)
}

func buildTable(headers []string, data [][]string) string {

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
