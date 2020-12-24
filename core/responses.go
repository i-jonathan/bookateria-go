package core

type response struct {
	Message string
}

var (
	// TwoHundred general response for http code 200
	TwoHundred	= response{Message: "OK"}
	// FourHundred general response for http code 400
	FourHundred = response{Message: "Invalid Request."}
	// FourOOne general response for http code 401
	FourOOne    = response{Message: "Access Denied."}
	// FourOFour general response for http code 404
	FourOFour   = response{Message: "Requested resource not found."}
	// FourTwoTwo general response for http code 422
	FourTwoTwo	= response{Message: "Your Request Could not be Processed."}
	// FiveHundred general response for http code 500
	FiveHundred = response{Message: "Server Error."}
)
