package entity

// HelloWorldRequest defines a request sent to the /helloworld endpoint
type HelloWorldRequest struct {
	Name string `json:"name"`
}

// HelloWorldResponse defines a response returned from the /helloworld endpoint
type HelloWorldResponse struct {
	Message string `json:"message"`
}
