package response

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	PerPage      int `json:"per_page"`
	LastPage     int `json:"last_page"`
	TotalResults int `json:"total_results"`
}
