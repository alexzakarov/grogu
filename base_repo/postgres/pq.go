package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/alexzakarov/grogu/config"
	"github.com/alexzakarov/grogu/utils"
	"strings"
)

type PQBaseRepo[C, U, G any] struct {
	db              *sql.DB
	PrimaryKey      string   `json:"primary_key"`
	Schema          string   `json:"schema"`
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

func NewPQBaseRepo[C, U, G any](config config.PQPostgresConfig) IBaseRepo[C, U, G] {
	var createStruc C
	var updateStruc U
	var getStruc G
	var createJsons []string
	var updateJsons []string
	var getJsons []string
	createReplacer := ""
	updateReplacer := ""
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
			createReplacer += fmt.Sprintf(`$%d,`, i+1)
		} else {
			createReplacer += fmt.Sprintf(`$%d`, i+1)
		}
	}
	for i, field := range updateJsons {
		if i < len(updateJsons)-1 {
			updateReplacer += fmt.Sprintf(`%s=$%d,`, field, i+1)
		} else {
			updateReplacer += fmt.Sprintf(`%s=$%d`, field, i+1)
		}
	}

	strCreateFields := strings.Join(createJsons, ",")
	strUpdateFields := strings.Join(updateJsons, ",")
	strGetFields := strings.Join(getJsons, ",")
	return &PQBaseRepo[C, U, G]{
		db:              config.Db,
		Schema:          config.Schema,
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

func (b *PQBaseRepo[C, U, G]) Create(dat C, success func(id int64), failure func(record int64)) {
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

	query := fmt.Sprintf(`INSERT INTO %s.%s (%s) VALUES (%s) RETURNING %s`, b.Schema, b.Table, b.strCreateFields, b.createReplacer, b.PrimaryKey)
	errDb = b.db.QueryRow(query, entity...).Scan(&data)
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

func (b *PQBaseRepo[C, U, G]) Update(entity_id int64, dat U, success func(), failure func(record int64)) {
	var mapEntity []interface{}
	var errParse error
	var errDb error
	var cmd sql.Result
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
		statusClause = fmt.Sprintf(`AND %s=%v`, b.statusName, ConvertStatus(1, b.statusType))
	}

	query := fmt.Sprintf(`UPDATE %s.%s SET %s WHERE %s=$%d %s`, b.Schema, b.Table, b.updateReplacer, b.PrimaryKey, len(entity), statusClause)
	cmd, errDb = b.db.Exec(query, entity...)
	affected, errRows := cmd.RowsAffected()
	if errRows != nil {
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

func (b *PQBaseRepo[C, U, G]) GetOne(entity_id int64, success func(data G), failure func(record int64), sub_queries ...SubQuery) {
	var bytes []byte
	var errDb error
	var entity G
	var qsReplaced []string
	var statusClause string

	for _, q := range sub_queries {
		var injected string
		replaced := strings.ReplaceAll(q.Query, "[primary_key]", fmt.Sprintf("ent.%s", b.PrimaryKey))
		if q.IsSingle {
			injected = fmt.Sprintf(`,(SELECT TO_JSON(ENTITY) FROM (%s) ENTITY) as %s`, replaced, q.Alias)
		} else {
			injected = fmt.Sprintf(`,(SELECT JSON_AGG(ENTITY) FROM (%s) ENTITY) as %s`, replaced, q.Alias)
		}
		qsReplaced = append(qsReplaced, injected)
	}

	subQs := strings.Join(qsReplaced, " ")

	if b.softDeletable {
		statusClause = fmt.Sprintf(`AND %s=%v`, b.statusName, ConvertStatus(1, b.statusType))
	}
	query := fmt.Sprintf(`SELECT TO_JSON(ENTITY) FROM (SELECT %s %s FROM %s.%s ent WHERE %s=$1 %s) ENTITY`, b.strGetFields, subQs, b.Schema, b.Table, b.PrimaryKey, statusClause)
	errDb = b.db.QueryRow(query, entity_id).Scan(&bytes)
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

func (b *PQBaseRepo[C, U, G]) DeleteOne(entity_id int64, success func(), failure func(record int64)) {
	var query string
	var cmd sql.Result
	var errDb error
	var statusClause string
	var statusUpdateClause string

	if b.softDeletable {
		statusClause = fmt.Sprintf(`AND %s=%v`, b.statusName, ConvertStatus(1, b.statusType))
		statusUpdateClause = fmt.Sprintf(`%s=%v`, b.statusName, ConvertStatus(2, b.statusType))
		query = fmt.Sprintf(`UPDATE %s.%s SET %s WHERE %s=$1 %s`, b.Schema, b.Table, statusUpdateClause, b.PrimaryKey, statusClause)
	} else {
		query = fmt.Sprintf(`DELETE FROM %s.%s WHERE %s=$1`, b.Schema, b.Table, b.PrimaryKey)
	}

	cmd, errDb = b.db.Exec(query, entity_id)
	affected, errRows := cmd.RowsAffected()
	if errRows != nil {
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

func (b *PQBaseRepo[C, U, G]) ChangeStatus(entity_id, status int64, success func(), failure func(record int64)) {
	var query string
	var cmd sql.Result
	var errDb error

	query = fmt.Sprintf(`UPDATE %s.%s SET %s=$1 WHERE %s=$2`, b.Schema, b.Table, b.statusName, b.PrimaryKey)
	cmd, errDb = b.db.Exec(query, ConvertStatus(status, b.statusType), entity_id)
	affected, errRows := cmd.RowsAffected()
	if errRows != nil {
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
