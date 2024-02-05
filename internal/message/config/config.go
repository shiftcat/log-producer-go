/*
	프로그램 실행 환경 설정 파일

	@author: yhan.lee shiftcats@gmail.com
    @date 2023-11-01
	@version: 0.1.0
*/

package config

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type KafkaConfig struct {
	BootstrapServers string                       `yaml:"bootstrap-servers"`
	TopicName        string                       `yaml:"topic-name"`
	Partitions       int32                        `yaml:"partitions"`
	Properties       map[string]kafka.ConfigValue `yaml:"properties"`
}

type DatabaseConfig struct {
	DriverName string `yaml:"driver-name"`
	DataSource string `yaml:"data-source"`
}

type Config struct {
	DataDir        string         `yaml:"data-dir"`
	NumberOfWorker int            `yaml:"number-of-worker"`
	ProfilePort    int            `yaml:"profile-port"`
	Kafka          KafkaConfig    `yaml:"kafka"`
	Database       DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) *Config {
	filename, _ := filepath.Abs(path)
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("config file load fail : %s", path))
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return &config
}

func (cfg *KafkaConfig) KafkaConfigMap() kafka.ConfigMap {
	m := make(map[string]kafka.ConfigValue)
	m["bootstrap.servers"] = cfg.BootstrapServers
	for k, v := range cfg.Properties {
		m[k] = v
	}
	return m
}
