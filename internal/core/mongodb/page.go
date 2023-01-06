package mongodb

//Page page
type Page struct {
	PageInformation PageInformation `json:"page_information,omitempty"`
	Entities        interface{}     `json:"entities,omitempty"`
}

// PageInformation page information
type PageInformation struct {
	Page                  int64 `json:"page"`
	Size                  int64 `json:"size"`
	TotalNumberOfEntities int64 `json:"total_number_of_entities"`
	TotalNumberOfPages    int64 `json:"total_number_of_pages"`
}
