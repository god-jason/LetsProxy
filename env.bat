

go env -w GOPROXY=https://goproxy.cn,direct

go env -w GOPRIVATE=*.gitlab.com,*.gitee.com,git.zgwit.com/xxx

go env -w GOSUMDB=off
:: go env -w GOSUMDB="sum.golang.google.cn"