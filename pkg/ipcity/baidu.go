package ipcity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type BaiDuIPInfo struct {
	ExtendedLocation string `json:"ExtendedLocation"`
	OriginQuery      string `json:"OriginQuery"`
	AppInfo          string `json:"appinfo"`
	DispType         int64  `json:"disp_type"`
	FetchKey         string `json:"fetchkey"`
	Location         string `json:"location"`
	OrigIp           string `json:"origip"`
	OrigIpQuery      string `json:"origipquery"`
}

type BaiDuIPRes struct {
	Status string         `json:"status"`
	Data   []*BaiDuIPInfo `json:"data"`
}

const baseUrl = "https://opendata.baidu.com/api.php?query=%s&co=&resource_id=6006&oe=utf8"

// GetGetLocationBaiDu 获取地址
func GetGetLocationBaiDu(address string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(baseUrl, address))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result BaiDuIPRes
	if err := json.Unmarshal(out, &result); err != nil {
		return "", err
	}
	if result.Status != "0" {
		return "", errors.New("get ip location error")
	}
	data := result.Data
	if len(data) == 0 {
		return "", nil
	}
	return data[0].Location, nil
}
