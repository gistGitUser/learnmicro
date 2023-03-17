package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "nics",
		Price: 1.00,
		SKU:   "abs-abc-def",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}

	p1 := &Product{
		Name:  "nics",
		Price: 1.00,
		SKU:   "abs-abcdef",
	}

	err2 := p1.Validate()

	if err2 == nil {
		t.Fatal("bad validation test")
	}
}
