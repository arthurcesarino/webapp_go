package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := new(Product)
	p.Name = "Arthur"
	p.Price = 1.00
	p.SKU = "asd-dsd-asda"

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
