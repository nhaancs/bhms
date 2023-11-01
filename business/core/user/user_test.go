package user_test

import (
	"fmt"
	"github.com/nhaancs/bhms/business/data/dbtest"
	"github.com/nhaancs/bhms/foundation/docker"
	"testing"
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
}
