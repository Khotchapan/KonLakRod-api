package googleCloud

type Books struct {
	ID     string `json:"id" firestore:"id"`
	Name   string `json:"name" firestore:"name"`
	Author string `json:"author" firestore:"author"`
}

type CreateBooksForm struct {
	Name   *string `json:"name"`
	Author *string `json:"author"`
}

func (f *CreateBooksForm) Fill(data *Books) *Books {
	if f.Name != nil {
		data.Name = *f.Name
	}
	if f.Author != nil {
		data.Author = *f.Author
	}

	return data
}

type UpdateBooksForm struct {
	ID     *string `json:"id" validate:"required"`
	Name   *string `json:"name"`
	Author *string `json:"author"`
}

func (f *UpdateBooksForm) Fill(data *Books) *Books {
	if f.ID != nil {
		data.ID = *f.ID
	}
	if f.Name != nil {
		data.Name = *f.Name
	}
	if f.Author != nil {
		data.Author = *f.Author
	}
	return data
}

type DeleteUsersForm struct {
	ID *string `json:"id" validate:"required"`
}
