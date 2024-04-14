package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexzakarov/grogu/config"
	"github.com/alexzakarov/grogu/database/ports"
	"github.com/alexzakarov/grogu/examples"
	"github.com/alexzakarov/grogu/utils"
	"strings"
)

type SQLXBaseRepo[C, U, G any] struct {
	ctx             context.Context
	db              ports.IBaseDb
	PrimaryKey      string   `json:"primary_key"`
	Table           string   `json:"table"`
	createFields    []string `json:"create_fields"`
	updateFields    []string `json:"update_fields"`
	getFields       []string `json:"get_fields"`
	strCreateFields string   `json:"str_create_fields"`
	strUpdateFields string   `json:"str_update_fields"`
	strGetFields    string   `json:"str_get_fields"`
	createReplacer  string   `json:"create_replacers"`
	updateReplacer  string   `json:"update_replacer"`
	softDeletable   bool     `json:"has_status"`
	statusName      string   `json:"status_name"`
	statusType      string   `json:"status_type"`
}

func NewSQLXBaseRepo[C, U any, G examples.UserResDto](config config.PostgresConfig) config.IBaseRepo[C, U, G] {
	var createStruc C
	var updateStruc U
	var getStruc G
	var createJsons []string
	var updateJsons []string
	var getJsons []string
	createReplacer := ""
	updateReplacer := ""
	strGetFields := ""
	var errParse error

	createJsons, _, errParse = utils.Convert(createStruc)
	if errParse != nil {
		println(errParse.Error())
	}

	updateJsons, _, errParse = utils.Convert(updateStruc)
	if errParse != nil {
		println(errParse.Error())
	}

	getJsons, _, errParse = utils.Convert(getStruc)
	if errParse != nil {
		println(errParse.Error())
	}

	for i, _ := range createJsons {
		if i < len(createJsons)-1 {
			createReplacer += `?,`
		} else {
			createReplacer += `?`
		}
	}
	for i, field := range updateJsons {
		if i < len(updateJsons)-1 {
			updateReplacer += fmt.Sprintf(`%s=?,`, field)
		} else {
			updateReplacer += fmt.Sprintf(`%s=?`, field)
		}
	}

	for i, field := range getJsons {
		if i < len(getJsons)-1 {
			strGetFields += fmt.Sprintf(`'%s',%s,`, field, field)
		} else {
			strGetFields += fmt.Sprintf(`'%s',%s`, field, field)
		}
	}

	strCreateFields := strings.Join(createJsons, ",")
	strUpdateFields := strings.Join(updateJsons, ",")
	return &SQLXBaseRepo[C, U, G]{
		db:              config.Db,
		Table:           config.Table,
		createFields:    createJsons,
		updateFields:    updateJsons,
		getFields:       getJsons,
		PrimaryKey:      config.PrimaryKey,
		strCreateFields: strCreateFields,
		strUpdateFields: strUpdateFields,
		strGetFields:    strGetFields,
		createReplacer:  createReplacer,
		updateReplacer:  updateReplacer,
		softDeletable:   config.SoftDeletable,
		statusName:      config.StatusName,
		statusType:      config.StatusType,
	}
}

func (b *SQLXBaseRepo[C, U, G]) Create(dat C, success func(id int64), failure func(record int64)) {
	var mapEntity []interface{}
	var errParse error
	var errDb error
	var data int64
	entity := []interface{}{}

	_, mapEntity, errParse = utils.Convert(dat)
	if errParse != nil {
		println(errParse.Error())
		failure(-1)
		return
	}
	entity = append(entity, mapEntity...)

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, b.Table, b.strCreateFields, b.createReplacer)
	data, errDb = b.db.Insert(query, entity...)
	if errDb != nil && utils.CheckStringIfContains(errDb.Error(), "duplicate key value") == false {
		println(errDb.Error())
		failure(-1)
		return
	} else if errDb != nil && utils.CheckStringIfContains(errDb.Error(), "duplicate key value") == true {
		failure(-2)
		return
	}
	success(data)
	return
}

func (b *SQLXBaseRepo[C, U, G]) Update(entity_id int64, dat U, success func(), failure func(record int64)) {
	var mapEntity []interface{}
	var errParse error
	var errDb error
	var entity []interface{}
	var statusClause string

	_, mapEntity, errParse = utils.Convert(dat)
	if errParse != nil {
		println(errParse.Error())
		failure(-1)
		return
	}
	entity = append(entity, mapEntity...)
	entity = append(entity, entity_id)

	if b.softDeletable {
		statusClause = fmt.Sprintf(`AND %s=%v`, b.statusName, utils.ConvertStatus(1, b.statusType))
	}

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE %s=? %s`, b.Table, b.updateReplacer, b.PrimaryKey, statusClause)
	_, errDb = b.db.Exec(query, entity...)
	if errDb != nil {
		println(errDb.Error())
		failure(-1)
		return
	}
	success()
	return
}

func (b *SQLXBaseRepo[C, U, G]) GetOne(entity_id int64, success func(data G), failure func(record int64), sub_queries ...config.SubQuery) {
	var bytes []byte
	var errDb error
	var entity G
	var qsReplaced []string
	var statusClause string

	for _, q := range sub_queries {
		var injected string
		replaced := strings.ReplaceAll(q.Query, "[primary_key]", fmt.Sprintf("ent.%s", b.PrimaryKey))
		if q.IsSingle {
			injected = fmt.Sprintf(`,(SELECT JSON_OBJECT(ENTITY) FROM (%s) ENTITY) as %s`, replaced, q.Alias)
		} else {
			injected = fmt.Sprintf(`,(SELECT JSON_ARRAYAGG(JSON_OBJECT(ENTITY)) FROM (%s) ENTITY) as %s`, replaced, q.Alias)
		}
		qsReplaced = append(qsReplaced, injected)
	}

	subQs := strings.Join(qsReplaced, " ")

	if b.softDeletable {
		statusClause = fmt.Sprintf(`AND %s=%v`, b.statusName, utils.ConvertStatus(1, b.statusType))
	}
	query := fmt.Sprintf(`SELECT JSON_OBJECT(%s) %s FROM %s ent WHERE %s=? %s`, b.strGetFields, subQs, b.Table, b.PrimaryKey, statusClause)
	bytes, errDb = b.db.Select(query, entity_id)
	if errDb != nil && utils.CheckStringIfContains(errDb.Error(), "no rows in result set") == false {
		println(errDb.Error())
		failure(-1)
		return
	} else if errDb != nil && utils.CheckStringIfContains(errDb.Error(), "no rows in result set") == true {
		failure(0)
		return
	}
	err := json.Unmarshal(bytes, &entity)
	if err != nil {
		println(err.Error())
		failure(-2)
	}
	success(entity)
	return
}

func (b *SQLXBaseRepo[C, U, G]) DeleteOne(entity_id int64, success func(), failure func(record int64)) {
	var query string
	var affected int64
	var errDb error
	var statusClause string
	var statusUpdateClause string

	if b.softDeletable {
		statusClause = fmt.Sprintf(`AND %s=%v`, b.statusName, utils.ConvertStatus(1, b.statusType))
		statusUpdateClause = fmt.Sprintf(`%s=%v`, b.statusName, utils.ConvertStatus(2, b.statusType))
		query = fmt.Sprintf(`UPDATE %s SET %s WHERE %s=? %s`, b.Table, statusUpdateClause, b.PrimaryKey, statusClause)
	} else {
		query = fmt.Sprintf(`DELETE FROM %s WHERE %s=?`, b.Table, b.PrimaryKey)
	}

	affected, errDb = b.db.Exec(query, entity_id)
	if errDb != nil {
		println(errDb.Error())
		failure(-1)
		return
	} else if affected == 0 {
		failure(0)
		return
	}
	success()
	return
}

func (b *SQLXBaseRepo[C, U, G]) ChangeStatus(entity_id, status int64, success func(), failure func(record int64)) {
	var query string
	var affected int64
	var errDb error

	query = fmt.Sprintf(`UPDATE %s SET %s=? WHERE %s=?`, b.Table, b.statusName, b.PrimaryKey)
	affected, errDb = b.db.Exec(query, utils.ConvertStatus(status, b.statusType), entity_id)
	if errDb != nil {
		println(errDb.Error())
		failure(-1)
		return
	} else if affected == 0 {
		failure(0)
		return
	}
	success()
	return
}
