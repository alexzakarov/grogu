

<div align="center">
<img alt="go-grogu" src="/assets/go-grogu.png" width="400" />
</div>

<div align="center">

# ![grogu](./assets/grogu-text.svg)

</div>

### Grogu is performance oriented ORM-Like library written in golang
 
## Initialization 

<p> The library uses struct as db fields. This means you do not need to specify the db fields more than once.</p>

```
repo := base_repo.NewBaseRepo[CreateUserDbModel, UpdateUserDbModel, UserResDto](ctx, db, "public", "users", "user_id", false, "", "")
```

## Create Example
```
user := CreateUserReqDto{
	UserTitle: "Test User",
}
meta := user.ToDbModel("This user has admin role")

repo.Create(meta, func(id int64) {
	record = 1
}, func(rec int64) {
	// negative rec refers to db errors
	record = rec
})
```

## Update Example
```
user := UpdateUserReqDto{
	UserTitle: "Test User",
}
meta := user.ToDbModel("This user has admin role")

repo.Update(userId, meta, func() {
	record = 1
}, func(rec int64) {
	// negative rec refers to db errors
	record = rec
})
```

## GetOne Example
```
userId := 1

repo.GetOne(userId, func(user UserResDto) {
	record = 1
	resData = user
}, func(rec int64) {
	// negative rec refers to db errors
	record = rec
})
```

## DeleteOne Example
```
userId := 1

repo.DeleteOne(userId, func() {
	record = 1
}, func(rec int64) {
	// negative rec refers to db errors
	record = rec
})
```

## Note: This library is not production ready. This means, you should not use it in your scalable/heavy projects 

#
#

 <p align="center">May the Force be with You! </p>
