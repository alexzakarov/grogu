package config

type SubQuery struct {
	IsSingle bool   `json:"is_one"`
	Alias    string `json:"alias"`
	Query    string `json:"query"`
}

type IBaseRepo[C, U, G any] interface {
	Create(C, func(id int64), func(error))
	Update(int64, U, func(), func(error))
	GetOne(int64, func(G), func(error), ...SubQuery)
	DeleteOne(int64, func(), func(error))
	ChangeStatus(int64, int64, func(), func(error))
}
