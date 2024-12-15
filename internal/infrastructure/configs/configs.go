package configs

import (
	"github.com/spf13/viper"
)

var LocalizePath = ""

func Load(path, lcp, env, logPath string) (*Bootstrap, error) {
	// 读取环境配置
	if mode := env; mode != "" {
		Mode = EnvMode(mode)
	} else { // 默认「生产环境」
		Mode = Development
	}
	// 读取配置文件
	v := viper.New()
	v.SetConfigFile(filePathByMode(Mode, path))
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	//if config != nil {
	//	if err := v.Unmarshal(config); err != nil {
	//		return nil, err
	//	}
	//}
	var bc Bootstrap
	if err := v.Unmarshal(&bc); err != nil {
		return nil, err
	}
	if lcp == "" {
		lcp = path + "/localize"
	}
	LocalizePath = lcp
	bc.Log.OutPath = logPath
	return &bc, nil
}

func filePathByMode(mode EnvMode, path string) string {
	switch mode {
	case Development:
		path = path + "/config_dev.yaml"
	case Prerelease:
		path = path + "/config_pre.yaml"
	case Production:
		path = path + "/config_pro.yaml"
	}
	return path
}

func GwtI18RootPath() string {
	return LocalizePath
}
