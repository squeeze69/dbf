package dbf

import (
	"bytes"
	"reflect"
	"testing"
)

//ExampleUsage
/*
func ExampleUsage() {
	dbr, _ := dbf.NewReader(os.Stdin)
	// fmt.Printf("Mod date: %d-%d-%d\n", dbr.Year, dbr.Month, dbr.Day)
	fmt.Printf("Num records: %d\n", dbr.Length)
	// record is map[string]interface{}
}
*/

var testFile = bytes.NewReader([]byte{
	// Header:
	0x03, 0x6F, 0x07, 0x1A, 0x0D, 0x21, 0x00, 0x00, 0x81, 0x00, 0x55, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x4F, 0x42, 0x4A, 0x45, 0x43, 0x54, 0x49, 0x44, 0x00, 0x00, 0x00, 0x4E, 0x00, 0x00, 0x00, 0x00,
	0x0B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x4E, 0x61, 0x6D, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x43, 0x00, 0x00, 0x00, 0x00,
	0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x53, 0x68, 0x61, 0x70, 0x65, 0x5F, 0x4C, 0x65, 0x6E, 0x67, 0x00, 0x46, 0x00, 0x00, 0x00, 0x00,
	0x09, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x0D,

	// Start of record 1 - deleted flag:
	0x20,

	//	OBJECTID:
	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x31,

	// Name:
	0x41, 0x62, 0x62, 0x6F, 0x74, 0x73, 0x62, 0x75, 0x72, 0x79, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
	0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
	0x20, 0x20,

	// Shape_Leng:
	0x20, 0x30, 0x2E, 0x30, 0x35, 0x32, 0x34, 0x36, 0x37,
})

var reader *Reader

func init() {
	var err error
	reader, err = NewReader(testFile)
	if err != nil {
		panic(err)
	}
}

func TestModDate(t *testing.T) {
	y, m, d := reader.ModDate()
	if y != 2011 || m != 7 || d != 26 {
		t.Fatalf("wrong ModDate(): got %d-%d-%d, expected 2011-07-26\n", y, m, d) // also try t.Errorf()
	}
}

func TestFieldNames(t *testing.T) {
	actual := reader.FieldNames()
	expected := []string{"OBJECTID", "Name", "Shape_Leng"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("wrong FieldNames(): got %v, expected %v", actual, expected)
	}
}

func TestFieldTypes(t *testing.T) {
	var badFieldType = bytes.NewReader([]byte{
		// Header:
		0x03, 0x6F, 0x07, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x41, 0x00, 0x0C, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x4F, 0x42, 0x4A, 0x45, 0x43, 0x54, 0x49, 0x44, 0x00, 0x00, 0x00, 0x42, 0x00, 0x00, 0x00, 0x00,
		0x0B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x0D,
	})

	_, err := NewReader(badFieldType)
	expectedErr := "Sorry, dbf library doesn't recognize field type 'B'"
	if err.Error() != expectedErr {
		t.Fatalf("Expected error: %s\nbut got: %s", expectedErr, err)
	}
}

func TestOneRead(t *testing.T) {
	expected := Record{
		"OBJECTID":   1,
		"Name":       "Abbotsbury",
		"Shape_Leng": 0.052467,
	}
	actual, err := reader.Read(0)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Read(0) returned wrong result: got %#v, expected %#v", actual, expected)
	}
}

func TestConcurrentReads(t *testing.T) {
	go TestOneRead(t)
	go TestOneRead(t)
}
