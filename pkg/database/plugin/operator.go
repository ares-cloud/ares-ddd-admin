package plugin

import (
	"reflect"

	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"gorm.io/gorm"
)

const (
	Creator = "creator"
	Updater = "updater"
)

type OperatorPlugin struct{}

func (t *OperatorPlugin) Name() string {
	return "operator_plugin"
}

func NewOperatorPlugin() *TenantPlugin {
	return &TenantPlugin{}
}

func (t *OperatorPlugin) Initialize(db *gorm.DB) error {
	if err := db.Callback().Create().Before("gorm:create").Register("user_tracking:before_create", t.beforeCreate); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").Register("user_tracking:before_update", t.beforeUpdate); err != nil {
		return err
	}
	return nil
}

func (t *OperatorPlugin) beforeCreate(db *gorm.DB) {
	ctx := db.Statement.Context
	OperatorUserId := actx.GetUserId(ctx)
	if db.Statement.Schema != nil && OperatorUserId != "" {
		// 检查是否存在租户字段
		field := db.Statement.Schema.FieldsByDBName[Creator]
		if field == nil {
			// 表中没有租户字段，直接返回
			return
		}
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				rv := reflect.Indirect(db.Statement.ReflectValue.Index(i))
				field1 := db.Statement.Schema.FieldsByDBName[Creator]
				if field1 != nil {
					err := field1.Set(ctx, rv, OperatorUserId)
					if err != nil {
						err = db.Statement.AddError(err)
						if err != nil {
							return
						}
					}
				}

			}
		case reflect.Struct:
			field := db.Statement.Schema.FieldsByDBName[Creator]
			if field != nil {
				db.Statement.SetColumn(Creator, OperatorUserId)
			}
		default:

		}
	}
}

func (t *OperatorPlugin) beforeUpdate(db *gorm.DB) {
	ctx := db.Statement.Context
	OperatorUserId := actx.GetUserId(ctx)
	if db.Statement.Schema != nil && OperatorUserId != "" {
		// 检查是否存在租户字段
		field := db.Statement.Schema.FieldsByDBName[Updater]
		if field == nil {
			// 表中没有租户字段，直接返回
			return
		}
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				rv := reflect.Indirect(db.Statement.ReflectValue.Index(i))
				field1 := db.Statement.Schema.FieldsByDBName[Updater]
				if field1 != nil {
					err := field1.Set(ctx, rv, OperatorUserId)
					if err != nil {
						err = db.Statement.AddError(err)
						if err != nil {
							return
						}
					}
				}

			}
		case reflect.Struct:
			field := db.Statement.Schema.FieldsByDBName[Updater]
			if field != nil {
				db.Statement.SetColumn(Updater, OperatorUserId)
			}
		default:

		}
	}
}
