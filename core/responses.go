package core

type Response struct {
	Message	string
}

var (
	FourOFour = Response{Message: "Requested resource not found."}
	FourOOne  = Response{Message: "Access Denied."}
	FourHundred = Response{Message: "Invalid Request."}
	FiveHundred = Response{Message: "Your request couldn't be processed by the server."}
)
