package fileutils

import (
	"crypto/sha256"
	"encoding/hex"
	"exponie/pkg/utils"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Close(f *os.File) {
	err := f.Close()
	if err != nil {
		log := log.Err(err)
		if f != nil {
			log = log.Str("file", f.Name())
		}
		log.Msg("Failed to close body")
	}
}

func MkdirParent(file string) error {
	log.Info().Str("path", file).Msg("Creating parent folders.")
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func Create(file string) (*os.File, error) {
	if err := MkdirParent(file); err != nil {
		return nil, err
	}
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Save(file string, data []byte) error {
	f, err := Create(file)
	if err != nil {
		return err
	}
	defer Close(f)
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func Copy(source, dest string) (*string, error) {
	f, err := Create(dest)
	if err != nil {
		return nil, err
	}
	defer Close(f)
	r, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer Close(r)
	hash := sha256.New()
	r2 := io.TeeReader(r, hash)
	_, err = io.Copy(f, r2)
	if err != nil {
		return nil, err
	}
	return utils.Ptr(hex.EncodeToString(hash.Sum(nil))), nil
}

func Sanitize(key string) string {
	key = filepath.Clean(filepath.Base(key))
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, " ", "_")
	return key
}

var homeDirectory = ""
var configDirectory = ""

func GetHomeDir() string {
	if homeDirectory == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Panic().Err(err).Msg("Failed to get home directory path")
		}
		homeDirectory = home
	}
	return homeDirectory
}

func GetConfigDir() string {
	if configDirectory == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			log.Panic().Err(err).Msg("Failed to get home directory path")
		}
		configDirectory = configDir
	}
	return configDirectory
}

func JoinHomePath(paths ...string) string {
	return filepath.Join(GetHomeDir(), filepath.Join(paths...))
}

func JoinConfigPath(paths ...string) string {
	return filepath.Join(GetConfigDir(), "exponie", filepath.Join(paths...))
}
