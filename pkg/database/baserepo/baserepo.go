package baserepo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/ares-cloud/ares-ddd-admin/pkg/database"
	"github.com/ares-cloud/ares-ddd-admin/pkg/database/query"
	"golang.org/x/exp/constraints"
	"gorm.io/gorm"
)

// SupportedIDTypes 支持的 ID 类型
type SupportedIDTypes interface {
	constraints.Integer | ~string
}

type IModel interface {
	GetPrimaryKey() string // 获取主键
	TableName() string     //模型表名
}

// IBaseRepo [T IModel,I SupportedIDTypes] ， 基础数据层方法
type IBaseRepo[T IModel, I SupportedIDTypes] interface {
	FindById(ctx context.Context, id I) (*T, error)                   // 根据 id 获取模型
	FindByIds(ctx context.Context, ids []I) ([]*T, error)             // 根据 id 获取模型
	DelById(ctx context.Context, id I) error                          // 根据 id 删除
	DelByIds(ctx context.Context, ids []I) error                      // 根据 id 批量删除
	DelByIdUnScoped(ctx context.Context, id I) error                  // 根据 id 物理删除（可单个可批量）
	DelByIdsUnScoped(ctx context.Context, ids []I) error              // 根据 id 物理删除（可单个可批量）
	EditById(ctx context.Context, id I, data *T) error                // 上下文/id/需要更新的数据模型或者map
	Add(ctx context.Context, data *T) (*T, error)                     // 创建并返回模型
	BathAdd(ctx context.Context, data ...*T) error                    // 批量插入
	Count(ctx context.Context, qb *query.QueryBuilder) (int64, error) // 计数
	Find(ctx context.Context, qb *query.QueryBuilder) ([]*T, error)   // 查询
	Db(ctx context.Context) *gorm.DB                                  // 获取db
	GetDb() database.IDataBase                                        // 获取database
	GenStringId() string
	GenInt64Id() int64
}

// BaseRepo [T interface{}] ， 基础数据层方法
type BaseRepo[T IModel, I SupportedIDTypes] struct {
	Model T                  // 模型
	db    database.IDataBase // 数据库连接
}

func NewBaseRepo[T IModel, I SupportedIDTypes](db database.IDataBase, model T) *BaseRepo[T, I] {
	return &BaseRepo[T, I]{
		db:    db,
		Model: model,
	}
}

// FindById ， 根据 id 获取模型
// 参数：
//
//	ctx ： 上下文
//	id ： 模型 id
//
// 返回值：
//
//	*T ：desc
//	error ：desc
func (r *BaseRepo[T, I]) FindById(ctx context.Context, id I) (*T, error) {
	var res T
	resDb := r.db.DB(ctx).Model(r.Model)
	v := reflect.ValueOf(r.Model)
	if v.FieldByName("DeletedAt").IsValid() {
		resDb.Where("deleted_at = 0")
	}
	//根据id查询
	if err := resDb.Where(fmt.Sprintf("%s = ?", r.Model.GetPrimaryKey()), id).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

// FindByIds ， 根据 id 获取模型
// 参数：
//
//	ctx ： 上下文
//	id ： 模型 id
//
// 返回值：
//
//	*T ：desc
//	error ：desc
func (r *BaseRepo[T, I]) FindByIds(ctx context.Context, ids []I) ([]*T, error) {
	var res []*T
	resDb := r.db.DB(ctx).Model(r.Model)
	v := reflect.ValueOf(r.Model)
	if v.FieldByName("DeletedAt").IsValid() {
		resDb.Where("deleted_at = 0")
	}
	//根据id查询
	if err := resDb.Where(fmt.Sprintf("%s in (?)", r.Model.GetPrimaryKey()), ids).First(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// DelById ， 根据 id 删除
// 参数：
//
//	ctx ： 上下文
//	id ： 模型 id
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T, I]) DelById(ctx context.Context, id I) error {
	db := r.db.DB(ctx)
	v := reflect.ValueOf(r.Model)
	if v.FieldByName("DeletedAt").IsValid() {
		if r.Model.GetPrimaryKey() == "" {
			// todo 可以用反射
			return errors.New("base repo model pk is not defined")
		}
		db.Model(r.Model).Where(fmt.Sprintf("%v = ?", r.Model.GetPrimaryKey()), id).Update("deleted_at", time.Now().Unix())
	} else {
		db.Delete(&r.Model, id)
	}
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// DelByIds ， 根据 id 批量删除
// 参数：
//
//	ctx ： 上下文
//	ids ： 模型 id 数组
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T, I]) DelByIds(ctx context.Context, ids []I) error {

	db := r.db.DB(ctx)
	v := reflect.ValueOf(r.Model)
	if v.FieldByName("DeletedAt").IsValid() {
		if r.Model.GetPrimaryKey() == "" {
			return errors.New("base repo model pk is not defined")
		}
		db.Model(r.Model).Where(fmt.Sprintf("%v in ?", r.Model.GetPrimaryKey()), ids).Update("deleted_at", time.Now().Unix())
	} else {
		db.Delete(&r.Model).Where(fmt.Sprintf("%v in ?", r.Model.GetPrimaryKey()), ids)
	}
	if err := db.Error; err != nil {
		return err
	}

	return nil
}

// DelByIdUnScoped ， 根据 id 删除(物理删除)
// 参数：
//
//	ctx ： 上下文
//	id ： 模型 id
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T, I]) DelByIdUnScoped(ctx context.Context, id I) error {
	return r.db.DB(ctx).Unscoped().Delete(&r.Model, id).Error
}

// DelByIdsUnScoped ， 根据 ids 批量删除(物理删除)
// 参数：
//
//	ctx ： 上下文
//	ids ： 模型 id 数组
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T, I]) DelByIdsUnScoped(ctx context.Context, ids []I) error {
	db := r.db.DB(ctx).Unscoped()
	v := reflect.ValueOf(r.Model)
	if v.FieldByName("DeletedAt").IsValid() {
		if r.Model.GetPrimaryKey() == "" {
			return errors.New("base repo model pk is not defined")
		}
		db.Model(r.Model).Where(fmt.Sprintf("%v in ?", r.Model.GetPrimaryKey()), ids).Update("deleted_at", time.Now().Unix())
	} else {
		db.Delete(&r.Model).Where(fmt.Sprintf("%v in ?", r.Model.GetPrimaryKey()), ids)
	}
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

// EditById ， 根据 id 更新 模型
// 参数：
//
//	ctx ： desc
//	id ： desc
//	data ： desc
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T, I]) EditById(ctx context.Context, id I, data *T) error {
	if r.Model.GetPrimaryKey() == "" {
		return errors.New("base repo model pk is not defined")
	}
	db := r.db.DB(ctx).Model(r.Model)
	v := reflect.ValueOf(r.Model)
	if v.FieldByName("DeletedAt").IsValid() {
		db.Where("deleted_at = 0")
	}
	v = reflect.ValueOf(data)
	if v.Kind() != reflect.Map {
		updated := v.Elem().FieldByName("UpdatedAt")
		if updated.IsValid() {
			updated.SetInt(time.Now().Unix())
		}
	}
	if err := db.Where(fmt.Sprintf("%s = ?", r.Model.GetPrimaryKey()), id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Add ， 创建模型
// 参数：
//
//	ctx ： 上下文
//	data ： 模型数据
//
// 返回值：
//
//	error ：desc
func (r *BaseRepo[T, I]) Add(ctx context.Context, data *T) (*T, error) {
	v := reflect.ValueOf(data)
	created := v.Elem().FieldByName("CreatedAt")
	if created.IsValid() {
		created.SetInt(time.Now().Unix())
	}
	if err := r.db.DB(ctx).Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *BaseRepo[T, I]) BathAdd(ctx context.Context, data ...*T) error {
	return r.db.DB(ctx).Create(data).Error
}

// 计数
func (r *BaseRepo[T, I]) Count(ctx context.Context, qb *query.QueryBuilder) (int64, error) {
	var count int64
	db := r.db.DB(ctx).Model(r.Model)
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}
	return count, db.Count(&count).Error
}

// 查询
func (r *BaseRepo[T, I]) Find(ctx context.Context, qb *query.QueryBuilder) ([]*T, error) {
	var res []*T
	db := r.db.DB(ctx).Model(r.Model)
	if where, values := qb.BuildWhere(); where != "" {
		db = db.Where(where, values...)
	}
	if orderBy := qb.BuildOrderBy(); orderBy != "" {
		db = db.Order(orderBy)
	}

	if limit, values := qb.BuildLimit(); limit != "" {
		db = db.Offset(values[0]).Limit(values[1])
	}
	return res, db.Find(&res).Error
}

func (r *BaseRepo[T, I]) Db(ctx context.Context) *gorm.DB {
	return r.db.DB(ctx)
}

func (r *BaseRepo[T, I]) GetDb() database.IDataBase {
	return r.db
}

func (r *BaseRepo[T, I]) GenStringId() string {
	return r.db.GenStringId()
}
func (r *BaseRepo[T, I]) GenInt64Id() int64 {
	return r.db.GenInt64Id()
}
