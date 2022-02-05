package handlers

import (
	"net/http"

	"github.com/arthurcesarino/webapp_go/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of products
// responses:
//	201: noContent
//	404: errorResponse
//	501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle DELETE method")

	id := getProductID(r)

	err := data.DeleteProducts(id)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Unable to delete product", http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
