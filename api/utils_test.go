// Copyright 2013 gandalf authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/globocom/commandmocker"
	"github.com/globocom/config"
	"github.com/globocom/gandalf/db"
	"github.com/globocom/gandalf/fs"
	"github.com/globocom/gandalf/user"
	testingfs "github.com/globocom/tsuru/fs/testing"
	"labix.org/v2/mgo/bson"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	tmpdir string
	rfs    *testingfs.RecordingFs
}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	err := config.ReadConfigFile("../etc/gandalf.conf")
	c.Assert(err, IsNil)
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "gandalf_api_tests")
	db.Connect()
	s.tmpdir, err = commandmocker.Add("git", "")
	c.Assert(err, IsNil)
}

func (s *S) SetUpTest(c *C) {
	s.rfs = &testingfs.RecordingFs{}
	fs.Fsystem = s.rfs
}

func (s *S) TearDownTest(c *C) {
	fs.Fsystem = nil
}

func (s *S) TearDownSuite(c *C) {
	commandmocker.Remove(s.tmpdir)
	db.Session.DB.DropDatabase()
}

func (s *S) TestGetUserOr404(c *C) {
	u := user.User{Name: "umi"}
	err := db.Session.User().Insert(&u)
	c.Assert(err, IsNil)
	defer db.Session.User().Remove(bson.M{"_id": u.Name})
	rUser, err := getUserOr404("umi")
	c.Assert(err, IsNil)
	c.Assert(rUser.Name, Equals, "umi")
}

func (s *S) TestGetUserOr404ShouldReturn404WhenUserDoesntExists(c *C) {
	_, e := getUserOr404("umi")
	expected := "User umi not found"
	got := e.Error()
	c.Assert(got, Equals, expected)
}
