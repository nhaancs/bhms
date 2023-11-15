package user_test

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/nhaancs/bhms/business/core/user"
	"github.com/nhaancs/bhms/business/data/dbtest"
	"github.com/nhaancs/bhms/foundation/docker"
	"runtime/debug"
	"testing"
	"time"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}

func Test_User(t *testing.T) {
	t.Run("crud", crud)
}

// =============================================================================

func crud(t *testing.T) {
	test := dbtest.NewTest(t, c)
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		test.Teardown()
	}()

	api := test.CoreAPIs

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// -------------------------------------------------------------------------

	newUsr := user.NewUser{
		FirstName: "Em",
		LastName:  "Teo",
		Phone:     "0984250068",
		Status:    user.StatusCreated,
		Roles:     []user.Role{user.RoleUser},
	}
	usr, err := api.User.Create(ctx, newUsr)
	if err != nil {
		t.Fatalf("Should be able to create user: %s.", err)
	}

	saved, err := api.User.QueryByID(ctx, usr.ID)
	if err != nil {
		t.Fatalf("Should be able to retrieve user by ID: %s.", err)
	}

	if !usr.CreatedAt.Equal(saved.CreatedAt) {
		t.Logf("got: %v", saved.CreatedAt)
		t.Logf("exp: %v", usr.CreatedAt)
		t.Logf("dif: %v", saved.CreatedAt.Sub(usr.CreatedAt))
		t.Errorf("Should get back the same created at")
	}

	if !usr.UpdatedAt.Equal(saved.UpdatedAt) {
		t.Logf("got: %v", saved.UpdatedAt)
		t.Logf("exp: %v", usr.UpdatedAt)
		t.Logf("dif: %v", saved.UpdatedAt.Sub(usr.UpdatedAt))
		t.Errorf("Should get back the same updated at")
	}

	if !usr.Status.Equal(saved.Status) {
		t.Logf("got: %v", saved.Status.Name())
		t.Logf("exp: %v", usr.Status.Name())
		t.Errorf("Should get back the same status")
	}

	if len(usr.Roles) != len(saved.Roles) {
		t.Logf("got: %v", len(saved.Roles))
		t.Logf("exp: %v", len(usr.Roles))
		t.Errorf("Should get back the same quantity roles")
	}

	for i := 0; i < len(usr.Roles); i++ {
		if !usr.Roles[i].Equal(saved.Roles[i]) {
			t.Logf("got: %v", len(saved.Roles[i].Name()))
			t.Logf("exp: %v", len(usr.Roles[i].Name()))
			t.Errorf("Should get back the same role at index %d", i)
		}
	}

	if diff := cmp.Diff(usr, saved, cmp.AllowUnexported(user.Role{}, user.Status{})); diff != "" {
		t.Fatalf("Should get back the same user. diff:\n%s", diff)
	}

	// -------------------------------------------------------------------------

	upd := user.UpdateUser{
		FirstName: dbtest.StringPointer("Nhan"),
		LastName:  dbtest.StringPointer("Nguyen"),
	}

	if _, err := api.User.Update(ctx, usr, upd); err != nil {
		t.Fatalf("Should be able to update user : %s.", err)
	}

	saved, err = api.User.QueryByPhone(ctx, usr.Phone)
	if err != nil {
		t.Fatalf("Should be able to retrieve user by Phone : %s.", err)
	}

	diff := usr.UpdatedAt.Sub(saved.UpdatedAt)
	if diff > 0 {
		t.Errorf("Should have a larger UpdatedAt : sav %v, usr %v, dif %v", saved.UpdatedAt, usr.UpdatedAt, diff)
	}

	if saved.FirstName != *upd.FirstName {
		t.Logf("got: %v", saved.FirstName)
		t.Logf("exp: %v", *upd.FirstName)
		t.Errorf("Should be able to see updates to FirstName")
	}

	if saved.LastName != *upd.LastName {
		t.Logf("got: %v", saved.LastName)
		t.Logf("exp: %v", *upd.LastName)
		t.Errorf("Should be able to see updates to LastName")
	}
}
