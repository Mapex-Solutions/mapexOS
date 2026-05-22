module github.com/Mapex-Solutions/MapexOS/contracts

go 1.25.3

require (
	github.com/Mapex-Solutions/mapexGoKit/infrastructure v0.0.0
	github.com/Mapex-Solutions/mapexGoKit/utils v0.0.0
)

require (
	github.com/Mapex-Solutions/mapexGoKit/microservices v0.0.0 // indirect
	github.com/klauspost/compress v1.18.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver/v2 v2.5.0 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
)

replace github.com/Mapex-Solutions/mapexGoKit/infrastructure => ../../../../mapexGoKit/infrastructure

replace github.com/Mapex-Solutions/mapexGoKit/microservices => ../../../../mapexGoKit/microservices

replace github.com/Mapex-Solutions/mapexGoKit/utils => ../../../../mapexGoKit/utils
