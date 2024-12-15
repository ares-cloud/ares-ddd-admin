package plugin

import (
	"context"
	"reflect"

	"github.com/ares-cloud/ares-ddd-admin/pkg/actx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TenantPlugin struct{}

func (t *TenantPlugin) Name() string {
	return "tenant_plugin"
}

func NewTenantPlugin() *TenantPlugin {
	return &TenantPlugin{}
}

func (t *TenantPlugin) Initialize(db *gorm.DB) error {
	if err := db.Callback().Query().Before("gorm:query").Register("tenant_id:before_query", t.beforeQuery); err != nil {
		return err
	}
	if err := db.Callback().Create().Before("gorm:create").Register("tenant_id:before_create", t.beforeCarte); err != nil {
		return err
	}
	return nil
}

// 创建前
func (t *TenantPlugin) beforeCarte(db *gorm.DB) {
	ctx := db.Statement.Context
	tenantID := actx.GetTenantId(ctx)
	if TenantIDNotNil(tenantID) && !IsIgnoreTenant(ctx) {
		if db.Statement.Schema != nil {
			// 检查是否存在租户字段
			field := db.Statement.Schema.FieldsByDBName[actx.KeyTenantId]
			if field == nil {
				// 表中没有租户字段，直接返回
				return
			}
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					rv := reflect.Indirect(db.Statement.ReflectValue.Index(i))
					field1 := db.Statement.Schema.FieldsByDBName[actx.KeyTenantId]
					if field1 != nil {
						err := field1.Set(ctx, rv, tenantID)
						if err != nil {
							err = db.Statement.AddError(err)
							if err != nil {
								return
							}
						}
					}

				}
			case reflect.Struct:
				field := db.Statement.Schema.FieldsByDBName[actx.KeyTenantId]
				if field != nil {
					db.Statement.SetColumn(actx.KeyTenantId, tenantID)
				}
			default:

			}
		}
	}
}

// 查询前
func (t *TenantPlugin) beforeQuery(db *gorm.DB) {
	// 一些业务逻辑，拿到 tenantID，可能从 context 中
	ctx := db.Statement.Context
	if !IsIgnoreTenant(ctx) {
		tenantID := actx.GetTenantId(ctx)
		if TenantIDNotNil(tenantID) && db.Statement.Schema != nil {
			// 检查是否存在租户字段
			field := db.Statement.Schema.FieldsByDBName[actx.KeyTenantId]
			if field == nil {
				// 表中没有租户字段，直接返回
				return
			}
			db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.Eq{Column: clause.Column{Table: db.Statement.Table, Name: actx.KeyTenantId}, Value: tenantID},
			}})
		}
		//租户为空的时候不加条件
		//else {
		//	db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{
		//		clause.Eq{Column: clause.Column{Table: db.Statement.Table, Name: hctx.KeyTenantId}, Value: ""},
		//	}})
		//}
	}
}

// TenantIDNotNil 租户id是否为空
func TenantIDNotNil(tenantID string) bool {
	return tenantID != "" && tenantID != "<nil>" && tenantID != "0"
}

// GetCtxTenantID 获取租户ID
func GetCtxTenantID(ctx context.Context) string {
	tenantID := actx.GetTenantId(ctx)
	if TenantIDNotNil(tenantID) {
		return tenantID
	}
	return ""
}

// IsIgnoreTenant 判断是否忽略租户
func IsIgnoreTenant(ctx context.Context) bool {
	return actx.IsIgnoreTenantId(ctx)
}

// AddTenantWhere 添加租户条件
func AddTenantWhere(ctx context.Context, db *gorm.DB, wStr string) *gorm.DB {
	tenantID := actx.GetTenantId(ctx)
	if TenantIDNotNil(tenantID) {
		db = db.Where(wStr, tenantID)
	} else {
		db = db.Where(wStr, "")
	}
	return db
}
