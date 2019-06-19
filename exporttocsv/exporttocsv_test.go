package exporttocsv

/*
func TestCreateCSV(t *testing.T) {
	// Create a file, check if it exists, and clean up after yourself.

	filename := "TestCSV"

	filename = CreateCSV(filename)
	f, err := os.Open(filename)
	if err != nil {
		t.Errorf("Error creating CSV. Check CreateCSV().\n")
	}
	f.Close()
	os.Remove(filename)
}
*/

/*
func TestExport(t *testing.T) {

	var err error
	filename := "TestCSV"
	record := []string{"item 1", "item2", "item3"}

	expected := 10
	for i := 0; i < expected; i++ {
		err = Export(filename, record)
	}
	if err != nil {
		t.Errorf("Error exporting record to CSV file.")
	}

	f, _ := os.Open("TestCSV.csv")
	defer os.Remove("TestCSV.csv")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	got := 0
	for scanner.Scan() {
		got++
	}
	if got != expected {
		t.Errorf("Expected: %d lines in the CSV file.\n", expected)
		t.Errorf("Got: %d lines.\n", got)
	}
}
*/
