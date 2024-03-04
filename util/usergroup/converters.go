// Package usergroup provides helpers for user and group
package usergroup

import (
	"fmt"
	"os/user"
	"strconv"
)

// UIDFromString function tries to retrieve user id from a user string 's'
// 's' may be an user id or a user name
func UIDFromString(s string) (uint32, error) {
	lookup, err := user.Lookup(s)
	if err != nil {
		lookup, err = user.LookupId(s)
		if err != nil {
			return 0, fmt.Errorf("unable to find user info for '%v'", s)
		}
	}
	var id int64
	id, err = strconv.ParseInt(lookup.Uid, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("unable to get userid info for '%v' (%v)", s, lookup.Uid)
	}
	userID := uint32(id)
	return userID, nil
}

// GIDFromString function tries to retrieve group id from a group string 's'
// 's' may be an group id or a group name
func GIDFromString(s string) (uint32, error) {
	lookup, err := user.LookupGroup(s)
	if err != nil {
		lookup, err = user.LookupGroupId(s)
		if err != nil {
			return 0, fmt.Errorf("unable to find group info for '%v'", s)
		}
	}
	var id int64
	id, err = strconv.ParseInt(lookup.Gid, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("unable to get groupid info for '%v' (%v)", s, lookup.Gid)
	}
	userID := uint32(id)
	return userID, nil
}
