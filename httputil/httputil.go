// Package httputil provides API response primitives and some useful functions
// to extend gin-gonic/gin's functionality.
package httputil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StdResponse contains basic elements of a standard API responses.
type StdResponse struct {
	Code     int                    `json:"code"`
	Message  string                 `json:"msg,omitempty"`
	Error    string                 `json:"err,omitempty"`
	Response map[string]interface{} `json:"resp"`
}

// NoField is an empty map[string]interface{}.
var NoField = map[string]interface{}{}

// Predefined standard API response.
var (
	//    Respond  | Class
	// ------------+-----------------
	//           0 | Unexpected Error
	//     1 -  99 | Common Error
	//   100 - 599 | HTTP Status
	//   600 - 699 | Set Up Error
	//   700 - 799 | Internal Error
	//   800 - 999 | API Error

	// UnexpectedErr is returned when any unknown errors occurred.
	UnexpectedErr = StdResponse{Code: 0, Message: "Something went wrong."}

	// HTTPSwitchProto is returned when the server has received the request headers and the client should proceed to send the request body.
	HTTPSwitchProto = StdResponse{Code: 100, Message: "Switching protocols."}

	// HTTPContinue is returned when the requester has asked the server to switch protocols and the server has agreed to do so.
	HTTPContinue = StdResponse{Code: 101, Message: "Continue."}

	// HTTPProcessing is returned when the server is requiring a long time to complete the request.
	HTTPProcessing = StdResponse{Code: 102, Message: "Processing."}

	// HTTPOkay is returned when everything are run as expected.
	HTTPOkay = StdResponse{Code: 200, Message: "OK."}

	// HTTPCreated is returned when the required resources are created.
	HTTPCreated = StdResponse{Code: 201, Message: "Created."}

	// HTTPAccepted is returned when the request is accepted.
	HTTPAccepted = StdResponse{Code: 202, Message: "Accepted."}

	// HTTPNoContent is returned when the request is processed, but nothing need to be return.
	HTTPNoContent = StdResponse{Code: 204, Message: "No content."}

	// HTTPNoResetContent is returned when the server requires the requester reset the document view.
	HTTPNoResetContent = StdResponse{Code: 205, Message: "Reset content."}

	// HTTPBadReq is returned when the client submit a bad request.
	HTTPBadReq = StdResponse{Code: 400, Message: "Bad request."}

	// HTTPUnauth is returned when the client is unaithorized to use this service.
	HTTPUnauth = StdResponse{Code: 401, Message: "Unauthorized."}

	// HTTPForbidden is returned when the client is forbidden to use this service.
	HTTPForbidden = StdResponse{Code: 403, Message: "Forbidden."}

	// HTTPNotFound is returned when the client is looking for an invaild path.
	HTTPNotFound = StdResponse{Code: 404, Message: "Not found."}

	// HTTPMethodNotAllow is returned when the client is using a invaild method.
	HTTPMethodNotAllow = StdResponse{Code: 405, Message: "Method not allowed."}

	// HTTPNotAcceptable is returned when the server cannot produce a response matching the list of acceptable values defined in the request's proactive content negotiation headers.
	HTTPNotAcceptable = StdResponse{Code: 406, Message: "Not acceptable."}

	// HTTPConflict is returned when a request conflict with current state of the server.
	HTTPConflict = StdResponse{Code: 409, Message: "Conflict."}
	
	// HTTPGone is returned when the target resource is no longer available.
	HTTPGone = StdResponse{Code: 410, Message: "Gone."}

	// HTTPUnsupportedMediaType is returned when the request payload is invaild.
	HTTPUnsupportedMediaType = StdResponse{Code: 415, Message: "Unsupported media type."}

	// HTTPTeapot is returned when the developer is drinking tea.
	HTTPTeapot = StdResponse{Code: 418, Message: "I'm a little teapot. Short and stout."}

	// HTTPUnprocessableEntity is returned when the content type of the request entity is correct, but the server was unable to process the contained instructions.
	HTTPUnprocessableEntity = StdResponse{Code: 422, Message: "Unprocessable entity."}

	// HTTPLocked is returned when the resource that is being accessed is locked.
	HTTPLocked = StdResponse{Code: 423, Message: "Locked."}

	// HTTPFailedDependency is returned when the request failed due to failure of a previous request.
	HTTPFailedDependency = StdResponse{Code: 424, Message: "Failed dependency."}

	// HTTPTooEarly is returned when the server is unwilling to risk processing a request that might be replayed..
	HTTPTooEarly = StdResponse{Code: 425, Message: "Too Early."}

	// HTTPInternalServerErr is returned when the server has encountered a situation it doesn't know how to handle.
	HTTPInternalServerErr = StdResponse{Code: 500, Message: "Internal server rrror."}

	// HTTPNotImplemented is returned when the client is looking for a function not implemented.
	HTTPNotImplemented = StdResponse{Code: 501, Message: "Not yet implemented, come back later."}

	// HTTPSvrUnavailable is returned when the server is unavailable to handle requests.
	HTTPSvrUnavailable = StdResponse{Code: 503, Message: "This service is temporarily unavailable, come back later."}

	// ExternalSvrErr is returned when errors occurred during calling external services.
	ExternalSvrErr = StdResponse{Code: 700, Message: "External service error."}

	// InvaildHeader is returned when required header fields are not provided by client.
	InvaildHeader = StdResponse{Code: 800, Message: "Required headers are invaild."}

	// InvaildParams is returned when the client provide invaild parameters.
	InvaildParams = StdResponse{Code: 801, Message: "Required parameters are invaild."}

	// NoRecordsFound is returned when no records are found in data source.
	NoRecordsFound = StdResponse{Code: 802, Message: "No records found in data source."}
)

// ToGinJSON return a map[string]interface{} with StdResponse and custom fields.
func (resp *StdResponse) ToGinJSON(err string, fields map[string]interface{}) (m map[string]interface{}) {
	resp.Error = err
	resp.Response = fields
	tempJSON, _ := json.Marshal(resp)
	json.Unmarshal(tempJSON, &m)
	return
}

// GinResp sets response in gin.Context body. This function can also log the full request and response to the logger if it is provided.
// Set the flag `showErr` to true if you want to return the `Error` field to client.
func GinResp(c *gin.Context, jsonObj map[string]interface{}, showErr bool, loggers ...*log.Logger) error {
	type ginLog struct {
		HTTPStatus int         `json:"http_status"`
		ClientIP   string      `json:"client_ip"`
		Path       string      `json:"path"`
		Params     string      `json:"params"`
		Header     http.Header `json:"header"`
		Body       string      `json:"body"`
		Response   interface{} `json:"response"`
	}

	// Log request and response if loggers are provided
	if len(loggers) > 0 {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		reqLog, err := json.Marshal(ginLog{
			HTTPStatus: http.StatusOK,
			ClientIP:   c.ClientIP(),
			Path:       c.FullPath(),
			Params:     c.Request.URL.Query().Encode(),
			Header:     c.Request.Header,
			Body:       string(body),
			Response:   jsonObj,
		})
		if err != nil {
			return err
		}
		for _, logger := range loggers {
			logger.Println(string(reqLog))
		}
	}

	// Set response to gin context
	if !showErr {
		delete(jsonObj, "err")
	}
	c.JSON(http.StatusOK, jsonObj)
	return nil
}

// GinRespWithAbort calls `c.Abort()` and then `GinResp()` internally.
func GinRespWithAbort(c *gin.Context, jsonObj map[string]interface{}, showErr bool, loggers ...*log.Logger) error {
	c.Abort()
	if err := GinResp(c, jsonObj, showErr, loggers...); err != nil {
		return err
	}
	return nil
}

// ShouldBindInfinite is similar with gin.ShouldBind, but it restores the request body after parsing it.
func ShouldBindInfinite(c *gin.Context, obj interface{}) error {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	if err := c.ShouldBind(obj); err != nil {
		return err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return nil
}

// BindHeaderInfinite is similar with gin.BindHeader, but it restores the request body after parsing it.
func BindHeaderInfinite(c *gin.Context, obj interface{}) error {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	if err := c.BindHeader(obj); err != nil {
		return err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	return nil
}
