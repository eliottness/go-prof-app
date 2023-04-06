module github.com/felixge/go-prof-app

go 1.16

require (
	github.com/DataDog/datadog-agent/pkg/remoteconfig/state v0.44.0-rc.4.0.20230328124545-8df2c0eb545f // indirect
	github.com/DataDog/datadog-go v4.8.2+incompatible
	github.com/DataDog/go-libddwaf v1.0.1-0.20230330154852-6f88cc26f53d // indirect
	github.com/ebitengine/purego v0.4.0-alpha.1.0.20230330151547-bfe9e39de4ce // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/jackc/pgx/v4 v4.14.0
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/nsrip-dd/cgotraceback v0.0.0-20220518170113-75f7f93d1852
	github.com/stretchr/testify v1.8.2 // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.39.0-alpha.1.0.20230406102041-5f67de397016
)

// replace gopkg.in/DataDog/dd-trace-go.v1 => /Users/felix.geisendoerfer/go/src/github.com/DataDog/dd-trace-go
