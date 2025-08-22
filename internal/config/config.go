package config

import "os"

// JWT
var JWTKey = []byte(os.Getenv("1234567890abcdef"))
