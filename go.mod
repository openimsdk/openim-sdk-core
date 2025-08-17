module github.com/openimsdk/openim-sdk-core/v3

go 1.20

require (
	github.com/coder/websocket v1.8.13
	github.com/golang/protobuf v1.5.4
	github.com/gorilla/websocket v1.5.0
	github.com/jinzhu/copier v0.4.0
	github.com/pkg/errors v0.9.1
	google.golang.org/protobuf v1.35.1 // indirect
	gorm.io/driver/sqlite v1.5.5
)

require (
	github.com/google/go-cmp v0.6.0
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/openimsdk/protocol v0.0.73-alpha.12
	github.com/openimsdk/tools v0.0.47-alpha.8
	github.com/patrickmn/go-cache v2.1.0+incompatible
	golang.org/x/image v0.24.0
	golang.org/x/sync v0.11.0
	gorm.io/gorm v1.25.10
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
)

replace (
	github.com/openimsdk/protocol => ./pkg/protocol
	github.com/openimsdk/tools => ./pkg/tools
	google.golang.org/protobuf => google.golang.org/protobuf v1.33.0
)
