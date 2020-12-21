package core

type Response struct {
	Message string
}

var (
	FourOFour   = Response{Message: "Requested resource not found."}
	FourOOne    = Response{Message: "Access Denied."}
	FourTwoTwo	= Response{Message: "Your Request Could not be Processed."}
	FourHundred = Response{Message: "Invalid Request."}
	FiveHundred = Response{Message: "Server Error."}
)
