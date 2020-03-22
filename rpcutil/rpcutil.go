package rpcutil

import (
	"encoding/json"
	"log"
)

// StdResponse contains basic elements of a standard exam service RPC response.
type StdResponse struct {
	Code int32
	Msg  string
	Err  string
	Resp map[string]interface{}
}

// NoField is an empty map[string]interface{}.
var NoField = map[string]interface{}{}

// Predefined standard RPC responses, inherit from standard HTTP API responses.
var (
	// UnexpectedErr is returned when any unknown errors occurred.
	UnexpectedErr = StdResponse{Code: 0, Msg: "Something went wrong."}

	// HTTPSwitchProto is returned when the server has received the request headers and the client should proceed to send the request body.
	HTTPSwitchProto = StdResponse{Code: 100, Msg: "Switching protocols."}

	// HTTPContinue is returned when the requester has asked the server to switch protocols and the server has agreed to do so.
	HTTPContinue = StdResponse{Code: 101, Msg: "Continue."}

	// HTTPProcessing is returned when the server is requiring a long time to complete the request.
	HTTPProcessing = StdResponse{Code: 102, Msg: "Processing."}

	// HTTPOkay is returned when everything are run as expected.
	HTTPOkay = StdResponse{Code: 200, Msg: "OK."}

	// HTTPCreated is returned when the required resources are created.
	HTTPCreated = StdResponse{Code: 201, Msg: "Created."}

	// HTTPAccepted is returned when the request is accepted.
	HTTPAccepted = StdResponse{Code: 202, Msg: "Accepted."}

	// HTTPNoContent is returned when the request is processed, but nothing need to be return.
	HTTPNoContent = StdResponse{Code: 204, Msg: "No content."}

	// HTTPNoResetContent is returned when the server requires the requester reset the document view.
	HTTPNoResetContent = StdResponse{Code: 205, Msg: "Reset content."}

	// HTTPBadReq is returned when the client submit a bad request.
	HTTPBadReq = StdResponse{Code: 400, Msg: "Bad request."}

	// HTTPUnauth is returned when the client is unaithorized to use this service.
	HTTPUnauth = StdResponse{Code: 401, Msg: "Unauthorized."}

	// HTTPForbidden is returned when the client is forbidden to use this service.
	HTTPForbidden = StdResponse{Code: 403, Msg: "Forbidden."}

	// HTTPNotFound is returned when the client is looking for an invaild path.
	HTTPNotFound = StdResponse{Code: 404, Msg: "Not found."}

	// HTTPMethodNotAllow is returned when the client is using a invaild method.
	HTTPMethodNotAllow = StdResponse{Code: 405, Msg: "Method not allowed."}

	// HTTPNotAcceptable is returned when the server cannot produce a response matching the list of acceptable values defined in the request's proactive content negotiation headers.
	HTTPNotAcceptable = StdResponse{Code: 406, Msg: "Not acceptable."}

	// HTTPConflict is returned when a request conflict with current state of the server.
	HTTPConflict = StdResponse{Code: 409, Msg: "Conflict."}

	// HTTPGone is returned when the target resource is no longer available.
	HTTPGone = StdResponse{Code: 410, Msg: "Gone."}

	// HTTPUnsupportedMediaType is returned when the request payload is invaild.
	HTTPUnsupportedMediaType = StdResponse{Code: 415, Msg: "Unsupported media type."}

	// HTTPTeapot is returned when the developer is drinking tea.
	HTTPTeapot = StdResponse{Code: 418, Msg: "I'm a little teapot. Short and stout."}

	// HTTPUnprocessableEntity is returned when the content type of the request entity is correct, but the server was unable to process the contained instructions.
	HTTPUnprocessableEntity = StdResponse{Code: 422, Msg: "Unprocessable entity."}

	// HTTPLocked is returned when the resource that is being accessed is locked.
	HTTPLocked = StdResponse{Code: 423, Msg: "Locked."}

	// HTTPFailedDependency is returned when the request failed due to failure of a previous request.
	HTTPFailedDependency = StdResponse{Code: 424, Msg: "Failed dependency."}

	// HTTPTooEarly is returned when the server is unwilling to risk processing a request that might be replayed..
	HTTPTooEarly = StdResponse{Code: 425, Msg: "Too Early."}

	// HTTPInternalServerErr is returned when the server has encountered a situation it doesn't know how to handle.
	HTTPInternalServerErr = StdResponse{Code: 500, Msg: "Internal server rrror."}

	// HTTPNotImplemented is returned when the client is looking for a function not implemented.
	HTTPNotImplemented = StdResponse{Code: 501, Msg: "Not yet implemented, come back later."}

	// HTTPSvrUnavailable is returned when the server is unavailable to handle requests.
	HTTPSvrUnavailable = StdResponse{Code: 503, Msg: "This service is temporarily unavailable, come back later."}

	// ExternalSvrErr is returned when errors occurred during calling external services.
	ExternalSvrErr = StdResponse{Code: 700, Msg: "External service error."}

	// InvaildHeader is returned when required header fields are not provided by client.
	InvaildHeader = StdResponse{Code: 800, Msg: "Required headers are invaild."}

	// InvaildParams is returned when the client provide invaild parameters.
	InvaildParams = StdResponse{Code: 801, Msg: "Required parameters are invaild."}

	// NoRecordsFound is returned when no records are found in data source.
	NoRecordsFound = StdResponse{Code: 802, Msg: "No records found in data source."}
)

// SetErr set the resp Error field.
func (resp *StdResponse) SetErr(err error) *StdResponse {
	resp.Err = err.Error()
	return resp
}

// SetResp set the resp Response field.
func (resp *StdResponse) SetResp(fields map[string]interface{}) *StdResponse {
	resp.Resp = fields
	return resp
}

// RPCResp logs the RPC request and returns the response to client.
func (resp *StdResponse) RPCResp(path string, req interface{}, loggers ...*log.Logger) (*StdResponse, error) {
	type rpcLog struct {
		Status   int32       `json:"status"`
		Path     string      `json:"path"`
		Params   interface{} `json:"params"`
		Response interface{} `json:"response"`
	}

	// Log request and response if loggers are provided
	if len(loggers) > 0 {
		reqLog, err := json.Marshal(rpcLog{
			Status:   resp.Code,
			Path:     path,
			Params:   req,
			Response: resp.Resp,
		})
		if err != nil {
			return resp, err
		}
		for i := range loggers {
			loggers[i].Println(string(reqLog))
		}
	}
	return resp, nil
}
