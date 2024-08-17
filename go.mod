module github.com/gdanko/enpass

go 1.22.2

replace gorm.io/driver/sqlite => github.com/open-olive/gorm-sqlcipher v1.1.4

require (
	github.com/BoredTape/gorm-sqlcipher v0.0.0-20210422223717-08e217b86307
	github.com/atotto/clipboard v0.1.4
	github.com/markkurossi/tabulate v0.0.0-20230223130100-d4965869b123
	github.com/miquella/ask v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.1
	golang.org/x/crypto v0.26.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/gorm v1.25.11
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
)
