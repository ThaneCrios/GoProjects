module gitlab.com/faemproject/backend/delivery/delivery.front

go 1.15

require (
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/casbin/casbin/v2 v2.23.0 // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-pg/pg v8.0.7+incompatible
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/korovkin/limiter v0.0.0-20190919045942-dac5a6b2a536
	github.com/labstack/echo/v4 v4.1.17
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/procfs v0.4.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	gitlab.com/faemproject/backend/core/shared v0.7.3
	go.uber.org/multierr v1.6.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace gitlab.com/faemproject/backend/core/shared => gitlab.com/faemproject/backend/core/shared.git v0.7.3
