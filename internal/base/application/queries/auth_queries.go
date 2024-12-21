package queries

type GetCaptchaQuery struct {
	Width  int64 `query:"width"`
	Height int64 `query:"height"`
}
