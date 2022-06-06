package main

import (
	"github.com/gin-gonic/gin"
	"go_demo/router"
)

func main() {
	//新增
	r := gin.Default()

	r = router.CollectRouter(r)

	panic(r.Run())

}
