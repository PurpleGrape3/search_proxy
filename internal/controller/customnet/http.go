package customnet

import (
	"bytes"
	"context"
	"net/http"
	"search_proxy/internal/model/objs"
	"search_proxy/internal/model/proxy"
	"search_proxy/internal/util/idgenerator"

	//	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type Net interface {
	StartNet(ip string, port string)
}

func NetFactory(netType string) Net {
	switch netType {
	case "http":
		return newCustomHttp()
	//TODO case "rpc":
	default:
		return newCustomHttp()
	}
}

type customHttp struct {
}

func newCustomHttp() *customHttp {
	return new(customHttp)
}

func (ch *customHttp) StartNet(ip string, port string) {
	/*router := gin.Default()
	equals
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())*/

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/add_doc", addDoc)
	router.GET("/del_doc", delDoc)
	router.GET("/doc_isdel", docIsDel)
	router.POST("/retrieve", retrieveDoc)
	//性能调试
	//pprof.Register(router)
	router.Run(ip + ":" + port)
}

func addDoc(ctx *gin.Context) {
	remoteIp := ctx.ClientIP()
	uri := ctx.Request.RequestURI
	bodyReader := ctx.Request.Body
	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyReader)
	body := buf.Bytes()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	retByte, err := proxy.AddDoc(newCtx, remoteIp, uri, body)

	if err != nil {
		var respData objs.RespData
		respData.Code = -1
		respData.Message = err.Error()
		ctx.JSON(http.StatusOK, respData)
	} else {
		ctx.Data(http.StatusOK, "application/json", retByte)
	}
}

func delDoc(ctx *gin.Context) {
	uri := ctx.Request.RequestURI
	remoteIp := ctx.ClientIP()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	retByte, err := proxy.DelDoc(newCtx, remoteIp, uri)

	if err != nil {
		var respData objs.RespData
		respData.Code = -1
		respData.Message = err.Error()
		ctx.JSON(http.StatusOK, respData)
	} else {
		ctx.Data(http.StatusOK, "application/json", retByte)
	}
}

func docIsDel(ctx *gin.Context) {
	uri := ctx.Request.RequestURI
	remoteIp := ctx.ClientIP()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	retByte, err := proxy.DocIsDel(newCtx, remoteIp, uri)

	if err != nil {
		var respData objs.RespData
		respData.Code = -1
		respData.Message = err.Error()
		ctx.JSON(http.StatusOK, respData)
	} else {
		ctx.Data(http.StatusOK, "application/json", retByte)
	}
}

func retrieveDoc(ctx *gin.Context) {
	remoteIp := ctx.ClientIP()
	uri := ctx.Request.RequestURI
	bodyReader := ctx.Request.Body
	buf := new(bytes.Buffer)
	buf.ReadFrom(bodyReader)
	body := buf.Bytes()
	trackid := uint64(idgenerator.Generate())
	newCtx := context.WithValue(ctx, "trackid", trackid)
	repl, errString := proxy.RetrieveDoc(newCtx, remoteIp, uri, body)

	if len(errString) != 0 {
		var respData objs.RespData
		respData.Code = -1
		respData.Message = errString
		respData.Result.Repl = repl
		ctx.JSON(http.StatusOK, respData)
	} else {
		var respData objs.RespData
		respData.Code = 0
		respData.Message = "ok"
		respData.Result.Repl = repl
		ctx.JSON(http.StatusOK, respData)
	}
}