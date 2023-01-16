package config

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Reader interface {
	RuntimeEnv() string
	Names() []string
	Bytes(name string) ([]byte, error)
}

type BinaryReader struct {
	names func() []string
	bytes func(string) ([]byte, error)
}

func NewBinaryReader(names func() []string, bytes func(string) ([]byte, error)) *BinaryReader {
	return &BinaryReader{
		names: names,
		bytes: bytes,
	}
}

func (b *BinaryReader) RuntimeEnv() string {
	return ""
}

func (b *BinaryReader) Names() []string {
	return b.names()
}

func (b *BinaryReader) Bytes(name string) ([]byte, error) {
	return b.bytes(name)
}

type RuntimeBinaryReader struct {
	*BinaryReader
	runtimeEnv string
}

func NewRuntimeBinaryReader(runtimeEnv string, names func() []string, bytes func(string) ([]byte, error)) *RuntimeBinaryReader {
	return &RuntimeBinaryReader{
		BinaryReader: &BinaryReader{
			names: names,
			bytes: bytes,
		},
		runtimeEnv: runtimeEnv,
	}
}

func (r *RuntimeBinaryReader) RuntimeEnv() string {
	return r.runtimeEnv
}

type Config struct {
	*viper.Viper
	runtimeEnv string
}

var c *Config = New()

func Default() *Config {
	return c
}

func New() *Config {
	c := &Config{
		Viper: viper.New(),
	}
	return c
}

func SetRuntimeEnv(s string) *Config {
	return c.SetRuntimeEnv(s)
}

func (c *Config) SetRuntimeEnv(s string) *Config {
	c.runtimeEnv = s
	return c
}

func Read(b Reader) *Config {
	return c.Read(b)
}

func (c *Config) Read(b Reader) *Config {
	v := viper.New()

	for _, name := range b.Names() {
		extWithPoint := filepath.Ext(name)
		if extWithPoint == "" {
			continue
		}
		ext := extWithPoint[1:]
		if !c.isExtValid(ext) {
			continue
		}
		v.SetConfigType(ext)
		data, err := b.Bytes(name)
		if err != nil {
			panic(err)
		}
		reader := bytes.NewReader(data)
		if err := v.MergeConfig(reader); err != nil {
			panic(err)
		}
	}

	// runtime sub
	var subSettings map[string]interface{}
	subViper := v.Sub(fmt.Sprintf("<%s>", c.runtimeEnv))
	if subViper != nil {
		subSettings = subViper.AllSettings()
	}

	// settings
	settings := v.AllSettings()
	for k := range settings {
		if k[0] == '<' && k[len(k)-1] == '>' {
			delete(settings, k)
		}
	}

	if b.RuntimeEnv() == "" || b.RuntimeEnv() == c.runtimeEnv {
		if err := c.MergeConfigMap(settings); err != nil {
			panic(err)
		}
		return c
	}

	// merge runtime config
	if c.runtimeEnv != "" {
		if err := v.MergeConfigMap(subSettings); err != nil {
			panic(err)
		}
	}

	return c
}

func (c *Config) isExtValid(ext string) bool {
	for _, e := range viper.SupportedExts {
		if ext == e {
			return true
		}
	}
	return false
}

func ReadBinary(names func() []string, bytes func(string) ([]byte, error)) *Config {
	return c.ReadBinary(names, bytes)
}

func (c *Config) ReadBinary(names func() []string, bytes func(string) ([]byte, error)) *Config {
	c.Read(NewBinaryReader(names, bytes))
	return c
}
