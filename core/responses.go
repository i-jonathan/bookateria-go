package core

type response struct {
	Message string
}

type ResponseStruct struct {
	Previous bool          `json:"previous"`
	Next     bool          `json:"next"`
	Page     int           `json:"page"`
	Count    int64         `json:"count"`
	Result   interface{} `json:"result"`
}


var (
	// TwoHundred general response for http code 200
	TwoHundred = response{Message: "OK"}
	// FourHundred general response for http code 400
	FourHundred = response{Message: "Invalid Request."}
	// FourOOne general response for http code 401
	FourOOne = response{Message: "Access Denied."}
	// FourOFour general response for http code 404
	FourOFour = response{Message: "Requested resource not found."}
	// FourONine response for http code 409
	FourONine = response{Message: "Conflict."}
	// FourTwoTwo general response for http code 422
	FourTwoTwo = response{Message: "Your Request Could not be Processed."}
	// FiveHundred general response for http code 500
	FiveHundred = response{Message: "Server Error."}
)
