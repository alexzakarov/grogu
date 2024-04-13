package postgres

type SubQuery struct {
	IsSingle bool   `json:"is_one"`
	Alias    string `json:"alias"`
	Query    string `json:"query"`
}

func ConvertStatus(status int64, status_type string) interface{} {
	var data interface{}
	if status_type == "bool" {
		switch status {
		case 1:
			data = true
		case 2:
			data = false
		}
	} else if status_type == "int" {
		data = status
	}
	return data
}

type IBaseRepo[C, U, G any] interface {
	Create(C, func(id int64), func(record int64))
	Update(int64, U, func(), func(int64))
	GetOne(int64, func(G), func(int64), ...SubQuery)
	DeleteOne(int64, func(), func(int64))
	ChangeStatus(int64, int64, func(), func(int64))
}
