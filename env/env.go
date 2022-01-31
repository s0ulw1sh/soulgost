package env

import (
	"os"
	"log"
	"strconv"
	"bufio"
	"strings"
	"path/filepath"
)

func GetEnvString(name, defval string) string {
	env := os.Getenv(name)
	
	if env == "" {
		return defval
	}

	return env
}

func GetEnvUint8(name string, defval uint8) uint8 {
	env := os.Getenv(name)
	
	if env == "" {
		return defval
	}

	val, err := strconv.ParseUint(env, 10, 8)

	if err != nil {
		log.Fatalf("invalid environment variable %s=%s, value must be a uint8", name, env)
		return defval
	}

	return uint8(val)
}

func GetEnvUint16(name string, defval uint16) uint16 {
	env := os.Getenv(name)
	
	if env == "" {
		return defval
	}

	val, err := strconv.ParseUint(env, 10, 16)

	if err != nil {
		log.Fatalf("invalid environment variable %s=%s, value must be a uint16", name, env)
		return defval
	}

	return uint16(val)
}

func GetEnvInt(name string, defval int) int {
	env := os.Getenv(name)
	
	if env == "" {
		return defval
	}

	val, err := strconv.ParseInt(env, 10, 32)

	if err != nil {
		log.Fatalf("invalid environment variable %s=%s, value must be a int", name, env)
		return defval
	}

	return int(val)
}

func GetEnvUint(name string, defval uint) uint {
	env := os.Getenv(name)
	
	if env == "" {
		return defval
	}

	val, err := strconv.ParseUint(env, 10, 32)

	if err != nil {
		log.Fatalf("invalid environment variable %s=%s, value must be a uint", name, env)
		return defval
	}

	return uint(val)
}

func GetEnvBool(name string, defval bool) bool {
	env := os.Getenv(name)
	
	if env == "" {
		return defval
	}

	b, err := strconv.ParseBool(env)

	if err != nil {
		log.Fatalf("invalid environment variable %s=%s, value must be a boolean", name, env)
		return defval
	}

	return b
}

func InitEnvFile() {
	wd, err := os.Getwd()

	if err != nil {
		return
	}

	file, err := os.Open(filepath.Join(wd, ".env"))

	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
			
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		sv := strings.Split(line, "=")

		if len(sv) > 1 {
			key := strings.TrimSpace(sv[0])
			val := strings.TrimSpace(strings.Join(sv[1:], "="))
			if key != "" {
				os.Setenv(key, val)
			}
		}
	}
}