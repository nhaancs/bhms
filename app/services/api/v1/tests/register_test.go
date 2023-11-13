package tests

import (
	"context"
	"github.com/nhaancs/bhms/business/core/user"
	"net/http"
	"os"
	"runtime/debug"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	v1 "github.com/nhaancs/bhms/app/services/api/v1"
	"github.com/nhaancs/bhms/app/services/api/v1/handlers/usergrp"
	"github.com/nhaancs/bhms/business/data/dbtest"
)

func createTests(t *testing.T, app appTest, sd seedData) {
	app.test(t, testCreate200(t, app, sd), "create200")
}

func testCreate200(t *testing.T, app appTest, sd seedData) []tableData {
	table := []tableData{
		{
			name:       "user",
			url:        "/v1/users/register",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
			model: &usergrp.AppRegister{
				FirstName: "Nhan",
				LastName:  "Nguyen",
				Phone:     "0984250068",
				Password:  "123456",
			},
			resp: &usergrp.AppUser{},
			expResp: &usergrp.AppUser{
				FirstName: "Nhan",
				LastName:  "Nguyen",
				Phone:     "0984250068",
				Status:    user.StatusCreated.Name(),
				Roles:     []string{"USER"},
			},
			cmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*usergrp.AppUser)
				expResp := y.(*usergrp.AppUser)

				if _, err := uuid.Parse(resp.ID); err != nil {
					return "bad uuid for ID"
				}

				if resp.CreatedAt == "" {
					return "missing date created"
				}

				if resp.UpdatedAt == "" {
					return "missing date updated"
				}

				expResp.ID = resp.ID
				expResp.CreatedAt = resp.CreatedAt
				expResp.UpdatedAt = resp.UpdatedAt

				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

// =============================================================================

func createSeed(ctx context.Context, api dbtest.CoreAPIs) (seedData, error) {
	return seedData{}, nil
}

// =============================================================================

func Test_Create(t *testing.T) {
	t.Parallel()

	dbTest := dbtest.NewTest(t, c)
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		dbTest.Teardown()
	}()

	h, err := v1.APIMux(v1.APIMuxConfig{
		Shutdown: make(chan os.Signal, 1),
		Log:      dbTest.Log,
		Auth:     dbTest.Auth,
		DB:       dbTest.DB,
	})
	if err != nil {
		t.Fatalf("APIMux error: %s", err)
	}

	app := appTest{
		Handler:    h,
		adminToken: dbTest.Token("0984250066", "gophers"),
		userToken:  dbTest.Token("0984250067", "gophers"),
	}

	// -------------------------------------------------------------------------

	t.Log("Seeding data ...")
	sd, err := createSeed(context.Background(), dbTest.CoreAPIs)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	createTests(t, app, sd)
}
