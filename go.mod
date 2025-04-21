module github.com/openimsdk/openim-sdk-core/v3

go 1.23.0

toolchain go1.24.0

require (
	github.com/golang/protobuf v1.5.4
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/copier v0.4.0
	github.com/pkg/errors v0.9.1
	google.golang.org/protobuf v1.35.1 // indirect
	gorm.io/driver/sqlite v1.5.5
	nhooyr.io/websocket v1.8.10
)

require golang.org/x/net v0.39.0 // indirect

require (
	github.com/google/go-cmp v0.6.0
	github.com/openimsdk/protocol v0.0.73-alpha.6
	github.com/openimsdk/tools v0.0.50-alpha.80
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/stretchr/testify v1.9.0
	golang.org/x/image v0.26.0
	golang.org/x/sync v0.13.0
	gorm.io/gorm v1.25.10
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/grpc v1.68.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace nhooyr.io/websocket => github.com/coder/websocket v1.8.10
