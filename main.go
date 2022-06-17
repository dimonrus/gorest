package gorest

// Meta information for pagination
type Meta struct {
	// Current page
	Page int `json:"page,omitempty"`
	// Limit
	Limit int `json:"limit,omitempty"`
	// Total count
	Total int `json:"total,omitempty"`
}

// HttpCodeMap Type for http code mapping
type HttpCodeMap map[string]int

// JsonResponse Json response struct
type JsonResponse struct {
	// HTTP code
	HttpCode int `json:"-"`
	// Message information
	Message interface{} `json:"message,omitempty"`
	// Body
	Data interface{} `json:"data,omitempty"`
	// Meta information
	Meta interface{} `json:"meta,omitempty"`
	// Error information
	Error interface{} `json:"error,omitempty"`
}
