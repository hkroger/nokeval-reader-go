package config

type SerialConfig struct {
	Device   string `yaml:"device"`
	Baud     uint   `yaml:"baud"`
	Bits     uint   `yaml:"bits"`
	StopBits uint   `yaml:"stopbits"`
}
