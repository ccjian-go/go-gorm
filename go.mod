module go-gorm

go 1.14

require (
    github.com/jinzhu/gorm v1.9.12
    my/tools v0.0.0
)

replace (
    my/tools => ./03/tools
)
