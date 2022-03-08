// Package middleware defines all HTTP middleware functions
// used to before handlers of routes of the pmd-dx-api.
package middleware

import (
	"context"
	"net/http"

	"github.com/janek64/pmd-dx-api/api/db"
	"github.com/janek64/pmd-dx-api/api/handler"
	"github.com/julienschmidt/httprouter"
)

// ResourceListParams checks for possible arguments of resource list queries, parses their
// values and stores them in a struct which is added to the context of the request.
func ResourceListParams(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Retrieve the parameters from the request
		queryParams := r.URL.Query()
		sort := queryParams.Get("sort")
		// Generate the ResourceListParams struct and add it to the context
		var params handler.ResourceListParams
		// Check if the value is one of the sort types
		if sort == db.IDAsc || sort == db.IDDesc || sort == db.NameAsc || sort == db.NameDesc {
			params.Sort.SortEnabled = true
			params.Sort.SortType = db.SortType(sort)
		} else {
			// Invalid ordering types are ignored instead of being answered with an error
			params.Sort.SortEnabled = false
		}
		ctx := context.WithValue(r.Context(), handler.ResourceListParamsKey, params)
		// Call the handler with the created context
		h(w, r.WithContext(ctx), ps)
	}
}
