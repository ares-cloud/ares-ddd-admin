package casbin

import "context"

type IPermissionsRepository interface {
	FindAllEnabled(context.Context) ([]*Role, error)
}
