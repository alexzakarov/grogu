package examples

type CreateUserReqDto struct {
	UserTitle string `json:"user_title"`
}

type CreateUserDbModel struct {
	MetaData  string `json:"meta_data"`
	UserTitle string `json:"user_title"`
}

func (m *CreateUserReqDto) ToDbModel(meta_data string) CreateUserDbModel {
	return CreateUserDbModel{
		MetaData:  meta_data,
		UserTitle: m.UserTitle,
	}
}

type UpdateUserReqDto struct {
	UserTitle string `json:"user_title"`
}

type UpdateUserDbModel struct {
	MetaData  string `json:"meta_data"`
	UserTitle string `json:"user_title"`
}

func (m *UpdateUserReqDto) ToDbModel(meta_data string) UpdateUserDbModel {
	return UpdateUserDbModel{
		MetaData:  meta_data,
		UserTitle: m.UserTitle,
	}
}

type UserResDto struct {
	UserId    int64  `json:"user_id"`
	MetaData  string `json:"meta_data"`
	UserTitle string `json:"user_title"`
}
