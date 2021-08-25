package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sarchimark/expense-tracker/app/handlers"
	"github.com/sarchimark/expense-tracker/business/tests"
	"github.com/sarchimark/expense-tracker/foundation/web"
)

type CategoryTests struct {
	app       http.Handler
	userToken string
}

func TestCategories(t *testing.T) {
	test := tests.NewIntegration(t,
		tests.DBContainer{
			Image: "postgres:latest",
			Port:  "5432",
			Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
		})

	t.Cleanup(test.Teardown)
	shutdown := make(chan os.Signal, 1)

	tests := CategoryTests{
		app:       handlers.SetupRoutes(test.DB, test.Log, shutdown),
		userToken: test.Token("sarchimark@example.com", "test123"),
	}

	t.Run("postCategory400", tests.postCategory400)

}

func (ct *CategoryTests) postCategory400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+ct.userToken)
	ct.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new category cant be created with an invalid document.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen using an incomplete category value.", testID)
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 400 for the response : %v", tests.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 400 for the response.", tests.Success, testID)

			//inspect response
			var got web.ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response to an error type : %v", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to unmarshal the response to an error type.", tests.Success, testID)

			//Define what we want to see.
			exp := web.ErrorResponse{
				Error: "field validation error",
				Fields: []web.FieldError{
					{Field: "title", Error: "title is a required field"},
					{Field: "description", Error: "description is a required field"},
				},
			}

			sorter := cmpopts.SortSlices(func(a, b web.FieldError) bool {

				return a.Field < b.Field
			})
			if diff := cmp.Diff(got, exp, sorter); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get the expected result. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get the expected result.", tests.Success, testID)

		}

	}

}
