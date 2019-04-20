package authz

import (
	"fmt"
	"strings"
)

// Permission is the authorization Permission which is formated like <service>.<resource>.<action>
// for example: iam.users.create
// iam.users.* can match all actions for iam.users permissions
type Permission struct {
	Service  PermissionField
	Resource PermissionField
	Action   PermissionField
}

// PermissionField is the permission value
type PermissionField string

// authz const
const (
	AllPermissionField = PermissionField("*")
)

// PermissionString is formated like <service>.<resource>.<action>
type PermissionString string

// PermissionsString is a string slice whose element is formated like a `PermissionString`
type PermissionsString []string

// Parse create a `*Permission` from a string
func (ps PermissionString) Parse() *Permission {
	split := strings.Split(string(ps), ".")
	if len(split) != 3 {
		return nil
	}

	return NewPermission(split[0], split[1], split[2])
}

// IsGranted to check if ps has permission of `needs`
func (ps PermissionString) IsGranted(needs *Permission) bool {
	if has := ps.Parse(); has != nil {
		return has.IsGranted(needs)
	}
	return false
}

// IsGranted to check if pss has permission of `needs`
func (pss PermissionsString) IsGranted(needs *Permission) bool {
	for _, ps := range pss {
		if has := PermissionString(ps).Parse(); has != nil {
			if has.IsGranted(needs) {
				return true
			}
		}
	}
	return false
}

func (pf PermissionField) isGranted(needs PermissionField) bool {
	if pf == AllPermissionField {
		return true
	}

	return pf == needs
}

// ID returns the permission identifier
func (p *Permission) ID() string {
	return fmt.Sprintf("%s.%s.%s", p.Service, p.Resource, p.Action)
}

func (p *Permission) String() string {
	return p.ID()
}

// IsGranted to check if p has permission of `needs`
func (p *Permission) IsGranted(needs *Permission) bool {
	return p.Service.isGranted(needs.Service) &&
		p.Resource.isGranted(needs.Resource) &&
		p.Action.isGranted(needs.Action)
}

// NewPermission create a new Permission
func NewPermission(service, resource, action string) *Permission {
	return &Permission{
		Service:  PermissionField(service),
		Resource: PermissionField(resource),
		Action:   PermissionField(action),
	}
}
