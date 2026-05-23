// Package export provides exporters for serialising drift detection results
// into portable file formats.
//
// Supported formats:
//
//   - CSV  — comma-separated values, suitable for spreadsheets and data pipelines.
//   - Markdown — GitHub-flavoured Markdown table, suitable for PR comments and reports.
//
// Example usage:
//
//	ex := export.New(os.Stdout)
//	if err := ex.WriteCSV(drifts); err != nil {
//		log.Fatal(err)
//	}
package export
