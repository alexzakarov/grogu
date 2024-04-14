package ports

type IBaseDb interface {
	Exec(string, ...interface{}) (int64, error)
	Insert(string, ...interface{}) (int64, error)
	Update(string, ...interface{}) (int64, error)
	Select(string, ...interface{}) ([]byte, error)
	Delete(string, ...interface{}) (int64, error)
}
