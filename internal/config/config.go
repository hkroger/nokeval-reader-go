package config

type MeasurementStorageConfig struct {
	DbType       string   `yaml:"type"`
	OverrideUrls []string `yaml:"override_url"`
	Secret       string   `yaml:"key"`
	ClientId     string   `yaml:"client_id"`
}

type Config struct {
	Serial             SerialConfig             `yaml:"serial"`
	MeasurementStorage MeasurementStorageConfig `yaml:"database"`
	FakeStorageMode    bool                     `yaml:"fake_storage_mode"`
	FakeSensorMode     bool                     `yaml:"fake_sensor_mode"`
	BufferFile         string                   `yaml:"buffer_file"`
}
