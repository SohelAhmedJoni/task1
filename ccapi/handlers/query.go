package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger-labs/cc-tools-demo/ccapi/chaincode"
	"github.com/hyperledger-labs/cc-tools-demo/ccapi/common"
)

func Query(c *gin.Context) {
	var args []byte
	var err error

	if c.Request.Method == "GET" {
		request := c.Query("@request")
		if request != "" {
			args, _ = base64.StdEncoding.DecodeString(request)
		}
	} else if c.Request.Method == "POST" {
		req := make(map[string]interface{})
		c.ShouldBind(&req)
		args, err = json.Marshal(req)
		if err != nil {
			common.Abort(c, http.StatusInternalServerError, err)
			return
		}
	}

	channelName := c.Param("channelName")
	chaincodeName := c.Param("chaincodeName")
	txName := c.Param("txname")

	argList := [][]byte{}
	if args != nil {
		argList = append(argList, args)
	}

	user := c.GetHeader("User")
	if user == "" {
		user = "Admin"
	}

	res, status, err := chaincode.Query(channelName, chaincodeName, txName, user, argList)
	if err != nil {
		common.Abort(c, status, err)
		return
	}

	var payload interface{}
	err = json.Unmarshal(res.Payload, &payload)
	if err != nil {
		common.Abort(c, http.StatusInternalServerError, err)
		return
	}

	common.Respond(c, payload, status, err)
}
