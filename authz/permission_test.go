package authz_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/westwin/gocommon/authz"
)

func TestPermissionMatch(t *testing.T) {
	needs := authz.NewPermission("iam", "users", "create")

	has := authz.PermissionString("iam.users.create")
	assert.True(t, has.IsGranted(needs), "%s should match %s", has, needs.ID())

	has = authz.PermissionString("iam.users.wrongaction")
	assert.False(t, has.IsGranted(needs), "%s should not match %s", has, needs.ID())
}

func TestPermissionWildcardMatch(t *testing.T) {
	needs := authz.NewPermission("iam", "users", "create")

	has := authz.PermissionString("iam.users.*")
	assert.True(t, has.IsGranted(needs), "%s should match %s", has, needs.ID())

	has = authz.PermissionString("*.*.*")
	assert.True(t, has.IsGranted(needs), "%s should match %s", has, needs.ID())

	has = authz.PermissionString("iam.*.create")
	assert.True(t, has.IsGranted(needs), "%s should match %s", has, needs.ID())

	has = authz.PermissionString("*.users.create")
	assert.True(t, has.IsGranted(needs), "%s should match %s", has, needs.ID())

	has = authz.PermissionString("*.users.wrongaction")
	assert.False(t, has.IsGranted(needs), "%s should not match %s", has, needs.ID())

	has = authz.PermissionString("iam.*.wrongaction")
	assert.False(t, has.IsGranted(needs), "%s should not match %s", has, needs.ID())

	has = authz.PermissionString("*.*.wrongaction")
	assert.False(t, has.IsGranted(needs), "%s should not match %s", has, needs.ID())
}

func TestParsePermission(t *testing.T) {
	service := "pn"
	resource := "push"
	action := "create"

	needsStr := authz.PermissionString(fmt.Sprintf("%s.%s.%s", service, resource, action))
	parsed := needsStr.Parse()

	assert.Equal(t, service, string(parsed.Service))
	assert.Equal(t, resource, string(parsed.Resource))
	assert.Equal(t, action, string(parsed.Action))
}

func TestNewPermission(t *testing.T) {
	service := "pn"
	resource := "push"
	action := "create"

	needs := authz.NewPermission(service, resource, action)
	assert.Equal(t, service, string(needs.Service))
	assert.Equal(t, resource, string(needs.Resource))
	assert.Equal(t, action, string(needs.Action))
}

func TestPermissionsStringMatch(t *testing.T) {
	needs := authz.NewPermission("iam", "users", "read")

	has := authz.PermissionsString{
		"iam.users.read",
	}
	assert.True(t, has.IsGranted(needs))

	has = authz.PermissionsString{
		"iam.users.create",
		"iam.users.read",
	}
	assert.True(t, has.IsGranted(needs))

	has = authz.PermissionsString{
		"*.*.*",
	}
	assert.True(t, has.IsGranted(needs))

	has = authz.PermissionsString{
		"*.*.*",
		"iam.users.create",
		"iam.users.read",
	}
	assert.True(t, has.IsGranted(needs))

	has = authz.PermissionsString{
		"iam.users.delete",
		"pn.push.create",
	}
	assert.False(t, has.IsGranted(needs))
}
