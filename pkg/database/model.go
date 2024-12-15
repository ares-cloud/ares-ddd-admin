package database

type BaseIntTime struct {
	CreatedAt int64 `json:"created_at" gorm:"column:created_at;not null;default:0;comment:创建时间"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at;not null;default:0;comment:更新时间"`
	DeletedAt int64 `json:"deleted_at" gorm:"column:deleted_at;not null;default:0;comment:删除时间"`
}

type BaseModel struct {
	BaseIntTime
	Creator string `json:"creator"  gorm:"column:creator;not null;default:'';comment:创建者"`
	Updater string `json:"updater"  gorm:"column:updater;not null;default:'';comment:更新人"`
}
