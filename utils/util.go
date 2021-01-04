package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JsonResponse struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Page    int64       `json:"page,omitempty"`
}

type Util struct {
}

type UtilInterface interface {
	SendErrorResponse(c *fiber.Ctx, msg string)
	SendBadRequestResponse(c *fiber.Ctx, msg string)
	IsZeroValue(x interface{}) bool
	ConvertStrArrToMongoObj(strArr []string) ([]primitive.ObjectID, error)
	ConvertInterfaceArrToStringArr(interfaceArr []interface{}) ([]string, error)
	//CommonRecover(fName string)
}

var UtilStruct UtilInterface

func init() {
	UtilStruct = &Util{}
}

func (u Util) SendErrorResponse(c *fiber.Ctx, msg string) {
	var resp []byte
	resp, _ = json.Marshal(JsonResponse{Success: false, Msg: msg})

	c.SendStatus(500)
	c.Send(resp)
}

func (u Util) SendBadRequestResponse(c *fiber.Ctx, msg string) {
	var resp []byte
	resp, _ = json.Marshal(JsonResponse{Success: false, Msg: msg})

	c.SendStatus(400)
	c.Send(resp)
}

//return true if zero value
func (u Util) IsZeroValue(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func (u Util) ConvertStrArrToMongoObj(strArr []string) ([]primitive.ObjectID, error) {
	var oids []primitive.ObjectID

	for _, sid := range strArr {
		soid, err := primitive.ObjectIDFromHex(sid)
		if err != nil {
			return make([]primitive.ObjectID, 0), err
		}
		oids = append(oids, soid)
	}

	return oids, nil
}

func (u Util) ConvertInterfaceArrToStringArr(interfaceArr []interface{}) ([]string, error) {
	var sids []string

	for _, id := range interfaceArr {
		sid := id.(string)
		sids = append(sids, sid)
	}

	return sids, nil
}

func CommonRecover(fName string) {
	if r := recover(); r != nil {
		fmt.Println(fName + ": stacktrace from panic: \n" + string(debug.Stack()))
	}
}
