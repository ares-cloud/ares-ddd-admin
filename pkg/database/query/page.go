package query

// Page 分页参数
type Page struct {
	Size    int `json:"size" query:"size"`       // 页码大小，最大100
	Current int `json:"current" query:"current"` // 页码，从1开始
}

func (p *Page) Offset() int {
	return (p.Current - 1) * p.Size
}

func (p *Page) Limit() int {
	return p.Size
}

func (p *Page) Fix() {
	if p.Current <= 0 {
		p.Current = 1
	}

	if p.Size <= 0 {
		p.Size = 10
	} else if p.Size > 500 {
		p.Size = 500
	}
}
