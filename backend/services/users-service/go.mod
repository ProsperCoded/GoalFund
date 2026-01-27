module github.com/gofund/users-service

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/gofund/shared v0.0.0
	gorm.io/gorm v1.25.5
)

replace github.com/gofund/shared => ../../shared
