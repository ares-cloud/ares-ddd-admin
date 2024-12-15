package models

type PageRes[T any] struct {
	Total int64 `json:"total"`
	List  []*T  `json:"list"`
}

func NewPageRes[T any](total int64, list []*T) *PageRes[T] {
	return &PageRes[T]{
		Total: total,
		List:  list,
	}
}

type BaseIntTime struct {
	CreatedAt int64 `json:"created_at"` //创建时间
	UpdatedAt int64 `json:"updated_at"` //更新时间
	DeletedAt int64 `json:"deleted_at"` //删除时间
}

type BaseModel struct {
	BaseIntTime
	Creator string `json:"creator"  gorm:"column:creator;not null;default:'';comment:创建者"`
	Updater string `json:"updater"  gorm:"column:updater;not null;default:'';comment:更新人"`
}
