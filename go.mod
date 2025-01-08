module github.com/openimsdk/openim-sdk-core/v3

go 1.22.7

toolchain go1.22.10

require (
	github.com/golang/protobuf v1.5.4
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/copier v0.4.0
	github.com/pkg/errors v0.9.1
	google.golang.org/protobuf v1.35.1
	gorm.io/driver/sqlite v1.5.5
	nhooyr.io/websocket v1.8.10
)

require golang.org/x/net v0.29.0 // indirect

require (
	github.com/google/go-cmp v0.6.0
	github.com/openimsdk/protocol v0.0.72-alpha.70
	github.com/openimsdk/tools v0.0.50-alpha.21
	github.com/patrickmn/go-cache v2.1.0+incompatible
	golang.org/x/image v0.15.0
	golang.org/x/sync v0.8.0
	google.golang.org/grpc v1.68.0
	gorm.io/gorm v1.25.10
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.9.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace nhooyr.io/websocket => github.com/coder/websocket v1.8.10
