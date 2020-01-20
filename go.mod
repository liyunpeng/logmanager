module logmanager

go 1.13

require (
	github.com/coreos/go-systemd/journal v0.0.0-00010101000000-000000000000 // indirect
	github.com/klauspost/cpuid v1.2.2 // indirect
	src/conf v0.0.0-00010101000000-000000000000 // indirect
	src/services v0.0.0-00010101000000-000000000000 // indirect
	src/web v0.0.0-00010101000000-000000000000
	src/web/controllers v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	github.com/coreos/go-systemd/journal => D:\gomodpah\go-systemd\journal
	src/conf => D:\goworkspace\logmanager\src\conf
	src/services => D:\goworkspace\logmanager\src\services
	src/web => D:\goworkspace\logmanager\src\web
	src/web/controllers => D:\goworkspace\logmanager\src\web\controllers
)
