package dto

type DateRangeRequest struct {
	Id        string `json:"id,omitempty"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
