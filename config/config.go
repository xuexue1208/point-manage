package config

import (
	"io/ioutil"
	"os"

	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v3"
)

type EnvParams struct {
	ListenAddr string `yaml:"listenaddr"`
	Mode       string `yaml:"mode"`
}

type DbInfos struct {
	DbType       string `yaml:"dbtype"`
	DbHost       string `yaml:"dbhost"`
	DbPort       int    `yaml:"dbport"`
	DbName       string `yaml:"dbname"`
	DbUser       string `yaml:"dbuser"`
	DbPwd        string `yaml:"dbpwd"`
	LogMode      bool   `yaml:"logmode"`
	MaxIdleConns int    `yaml:"maxidleconns"`
	MaxOpenConns int    `yaml:"maxopenconns"`
	MaxLifeTime  int    `yaml:"maxlifetime"`
}

type OSSConfig struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Host            string `yaml:"host"`
	CallbackUrl     string `yaml:"callback_url"`
	UploadDir       string `yaml:"upload_dir"`
	ExpireTime      int    `yaml:"expire_time"`
}
type DatabaseBusiness string //业务类型
type Config struct {
	EnvParams EnvParams                     `yaml:"env_params"`
	DbInfos   map[DatabaseBusiness]*DbInfos `yaml:"mysql"`
	OSSConfig map[string]*OSSConfig         `yaml:"oss"`
}

var Instance *Config //配置文件选项

func Initialize() {
	//读取配置文件
	if err := LoadConfig(); err != nil {
		logger.Fatal("%v", err)
	}
	logger.Info("config load  success")
}

func LoadConfig() error {
	vEnvironment := os.Getenv("env")
	defaultFile := "./conf/config.yaml"
	if vEnvironment == "local" {
		defaultFile = "./conf/config.yaml"
	}

	file, err := os.Open(defaultFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	c := &Config{}
	if err := yaml.Unmarshal(content, c); err != nil {
		return err
	}
	Instance = c

	return nil
}
