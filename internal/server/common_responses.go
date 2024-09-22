package server

var invalidPayloadResponse = HttpErrorResponse{
	Errors: []HttpError{
		{
			Code:    "Bad request",
			Title:   "Validation Error",
			Details: "Invalid payload format",
		},
	},
}

var internalServerErrorResponse = HttpErrorResponse{
	Errors: []HttpError{
		{
			Code:    "Internal Server Error",
			Title:   "Internal Server Error",
			Details: "Something went wrong",
		},
	},
}
