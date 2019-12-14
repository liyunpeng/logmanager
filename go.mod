module logmanager

go 1.13

require (
	github.com/Shopify/sarama v1.24.1 // indirect
	github.com/astaxie/beego v1.12.0 // indirect
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/klauspost/cpuid v1.2.2 // indirect
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	go.etcd.io/etcd v3.3.18+incompatible // indirect
	go.uber.org/zap v1.13.0 // indirect
	google.golang.org/grpc v1.25.1 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	src/conf v0.0.0-00010101000000-000000000000 // indirect
	src/services v0.0.0-00010101000000-000000000000 // indirect
	src/web v0.0.0-00010101000000-000000000000
	src/web/controllers v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	src/conf => D:\goworkspace\logmanager\src\conf
	src/services => D:\goworkspace\logmanager\src\services
	src/web => D:\goworkspace\logmanager\src\web
	src/web/controllers => D:\goworkspace\logmanager\src\web\controllers
)
