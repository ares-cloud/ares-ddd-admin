package models

type IntIdReq struct {
	Id int64 `query:"id" path:"id"` // 通用int类型Id(用于在path或param接受参数)
}

type StringIdReq struct {
	Id string `query:"id" path:"id"` // 通用string类型Id(用于在path或param接受参数)
}

// PageReq 分页参数
type PageReq struct {
	Size    int64 `query:"size,required"`    // 页码大小，最大100
	Current int64 `query:"current,required"` // 页码，从1开始
}
