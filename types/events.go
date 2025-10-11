package types

type Event struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type EventsResponse struct {
	Data  []Event `json:"data"`
	Items []Event `json:"items"`
}

type FetchResult struct {
	Event *Event
	Err   error
}
