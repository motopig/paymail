package paymail

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
)

type PayMail struct {
	Domain string
	// 检查paymail合法性
	CheckPayMail func(mail string) bool
	// 获取pubkey
	GetPubKey func(mail string) (string, error)
	// 验证pubkey paymail匹配
	VerifyPubMail func(pubkey string, mail string) (bool, error)
}

func New(p *PayMail) (*PayMail, error) {
	return p, nil
}

func Load(r *gin.Engine, payMailObj PayMail) {
	r.GET("/.well-known/bsvalias", func(c *gin.Context) {
		bsvalias(c, payMailObj)
	})
	api := r.Group("/api")
	api.GET("/v1/bsvalias/id/:paymail", func(c *gin.Context) {
		id(c, payMailObj)
	})
	api.GET("/v1/bsvalias/verify-pubkey/:paymail/:pubkey", func(c *gin.Context) {
		verify(c, payMailObj)
	})
	api.POST("/v1/bsvalias/address/:paymail", func(c *gin.Context) {
		address(c, payMailObj)
	})
}

func bsvalias(c *gin.Context, payMailObj PayMail) {
	data := ServiceDiscoveryResponse{
		Version: Version,
		Capabilities: Capabilities{
			SenderValidationUrl:     false,
			VerifyPublicKeyOwnerUrl: payMailObj.Domain + "/api/v1/bsvalias/verify-pubkey/{alias}@{domain.tld}/{pubkey}",
			PkiUrl:                  payMailObj.Domain + "/api/v1/bsvalias/id/{alias}@{domain.tld}",
			PaymentDestinationUrl:   payMailObj.Domain + "/api/v1/bsvalias/address/{alias}@{domain.tld}",
		},
	}
	c.JSON(http.StatusOK, data)
}

func notFound(paymail string) {
	data := NotFound{
		Code:    "not-found",
		Message: "Paymail not found: " + paymail,
	}
	c.JSON(http.StatusOK, data)
}

func id(c *gin.Context, payMailObj PayMail) {
	paymail := c.Param("paymail")
	if VerifyEmailFormat(paymail) == true {
		pubkey, err := payMailObj.GetPubKey(paymail)
		if err != nil {
			notFound(paymail)
			return
		}
		if pubkey != "" {
			data := PKIResponse{
				Version: Version,
				Handle:  paymail,
				PubKey:  pubkey,
			}
			c.JSON(http.StatusOK, data)
			return
		}
	}
	notFound(paymail)
}

func address(c *gin.Context, payMailObj PayMail) {
	//paymail := c.Param("paymail")
	// todo
	c.JSON(http.StatusOK, gin.H{"output": ""})
}

func verify(c *gin.Context, payMailObj PayMail) {
	paymail := c.Param("paymail")
	pubkey := c.Param("pubkey")
	data := VerifyResponse{
		Handle: paymail,
		Pubkey: pubkey,
		Match:  false,
	}
	if VerifyEmailFormat(paymail) == true {
		ok, err := payMailObj.VerifyPubMail(pubkey, paymail)
		if err != nil {
			log.Println(err.Error())
		}
		// todo
		if ok == true {
			data.Match = true
		}
		c.JSON(http.StatusOK, data)
		return
	}
	c.JSON(http.StatusOK, data)
}

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
