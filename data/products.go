package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this product
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this product
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`
}

// Products is a collection of Product
type Products []*Product

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func UpdateProducts(p *Product) error {
	pos := findIndexByProductID(p.ID)
	if pos == -1 {
		return ErrProductNotFound
	}

	productList[pos] = p
	return nil
}

// AddProduct adds a new product to the database
func AddProducts(p *Product) {
	p.ID = (productList[len(productList)-1].ID) + 1
	productList = append(productList, p)
}

// DeleteProduct deletes a product from the database
func DeleteProducts(id int) error {
	pos := findIndexByProductID(id)
	if pos == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:pos], productList[pos+1])
	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-absd-dfsdf
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "milky coffee",
		Price:       2.54,
		SKU:         "abc323",
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "strong milky coffee",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
