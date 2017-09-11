package mysql

import (
	"crm-search-task/lib/abstract"
	"crm-search-task/lib/utils"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"strings"
)

type DbConf struct {
	Host      string `ini:"host"`
	User      string `ini:"user"`
	Pass      string `ini:"pass"`
	Name      string `ini:"name"`
	MaxConns  int    `ini:"max_conns"`
	IdleConns int    `ini:"idle_conns"`
}

type Query struct {
	conf      *DbConf
	db        *sql.DB
	lastSql   string
	insertId  int
	affectRow int
}

const longQueryTime = 500 //ms

var log abstract.Logger

func New(conf *DbConf, logger abstract.Logger) (*Query, error) {
	q := &Query{conf: conf}
	err := q.createDb()
	if logger != nil {
		log = logger
	} else {
		log = &abstract.EmptyLog{}
	}
	return q, err
}

func (q *Query) createDb() (err error) {
	dbConf := q.conf
	if dbConf != nil {
		user := dbConf.User
		passwd := dbConf.Pass
		host := dbConf.Host
		database := dbConf.Name
		if q.db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", user, passwd, host, database)); err == nil {
			q.db.SetMaxIdleConns(q.conf.MaxConns)
			q.db.SetMaxOpenConns(q.conf.IdleConns)
		}

	}
	return
}

func (q *Query) ListMap(table, condition, fields string) (res []map[string]string, err error) {
	if rows, err := q.getRows(table, condition, fields); err == nil {
		defer rows.Close()
		res, err = rowsToMap(rows)
	}
	return
}

func (q *Query) QueryListMap(sql string) (res []map[string]string, err error) {
	if rows, err := q.Query(sql); err == nil {
		defer rows.Close()
		res, err = rowsToMap(rows)
	}
	return
}

func (q *Query) ListStruct(table, condition, fields string, list interface{}) (err error) {
	var dataList []map[string]string
	if dataList, err = q.ListMap(table, condition, fields); err == nil {
		return utils.ToStructList(dataList, list, "orm")
	}
	return
}

func (q *Query) QueryListStruct(sql string, list interface{}) (err error) {
	var dataList []map[string]string
	if dataList, err = q.QueryListMap(sql); err == nil {
		return utils.ToStructList(dataList, list, "orm")
	}
	return
}

func (q *Query) MapByUniqueField(table, condition, fields, uniqueField string) (res map[string]map[string]string, err error) {
	if rows, err := q.getRows(table, condition, fields); err == nil {
		defer rows.Close()
		var resList []map[string]string
		resList, err = rowsToMap(rows)
		res = make(map[string]map[string]string, len(resList))
		for _, item := range resList {
			if v, found := item[uniqueField]; found {
				res[v] = item
				delete(item, uniqueField)
			}
		}
	}
	return
}

func (q *Query) GetMap(table, condition, fields string) (res map[string]string, err error) {
	condition += " limit 1"
	var list []map[string]string
	if list, err = q.ListMap(table, condition, fields); err == nil && len(list) > 0 {
		res = list[0]
	}
	return
}

func (q *Query) QueryMap(sql string) (res map[string]string, err error) {
	var list []map[string]string
	if list, err = q.QueryListMap(sql); err == nil && len(list) > 0 {
		res = list[0]
	}
	return
}

func (q *Query) GetStruct(table, condition, fields string, stru interface{}) (err error) {
	var one map[string]string
	if one, err = q.GetMap(table, condition, fields); err == nil {
		err = utils.ToStruct(one, stru, "orm")
	}
	return
}

func (q *Query) QueryStruct(sql string, stru interface{}) (err error) {
	var one map[string]string
	if one, err = q.QueryMap(sql); err == nil {
		err = utils.ToStruct(one, stru, "orm")
	}
	return
}

func (q *Query) ListColumn(table, condition, fields string, col interface{}) (err error) {
	var rows *sql.Rows
	if rows, err = q.getRows(table, condition, fields); err == nil {
		defer rows.Close()
		return rowsToColumn(rows, col)
	}
	return
}

func (q *Query) QueryListColumn(sql string, col interface{}) (err error) {
	if rows, err := q.Query(sql); err == nil {
		defer rows.Close()
		return rowsToColumn(rows, col)
	} else {
		return err
	}
	return nil
}

func (q *Query) GetStringVal(table, condition, field string) (res string, err error) {
	var one map[string]string
	if one, err = q.GetMap(table, condition, field); err == nil {
		if len(one) == 1 {
			for _, v := range one {
				return v, nil
			}
		} else {
			err = errors.New("result err")
		}
	}
	return
}

func (q *Query) GetVal(table, condition, field string, res interface{}) (err error) {
	var rows *sql.Rows
	if rows, err = q.getRows(table, condition, field); err == nil {
		defer rows.Close()
		return rowsToVal(rows, res)
	}
	return nil
}

func (q *Query) QueryVal(sql string) (res string, err error) {
	var one map[string]string
	if one, err = q.QueryMap(sql); err == nil {
		if len(one) == 1 {
			for _, v := range one {
				return v, nil
			}
		} else {
			err = errors.New("result err")
		}
	}
	return
}

func (q *Query) getRows(table, condition, fields string) (*sql.Rows, error) {
	if len(fields) == 0 {
		fields = "*"
	}
	condition = strings.ToLower(strings.TrimSpace(condition))
	if len(condition) == 0 || strings.HasPrefix(condition, "order") || strings.HasPrefix(condition, "group") || strings.HasPrefix(condition, "limit") {
		condition = fmt.Sprintf("1=1 %s", condition)
	}
	querySql := fmt.Sprintf("select %s from `%s` where %s", fields, table, condition)
	return q.Query(querySql)
}

func (q *Query) Query(sql string, args ...interface{}) (rows *sql.Rows, err error) {
	cost := utils.FuncCost(func() {
		rows, err = q.db.Query(sql, args...)
		if err != nil {
			err = errors.New(fmt.Sprintf("%s | %s", err, sql))
		}
	})
	if cost > longQueryTime {
		log.Warnf("sql: %s, cost: %dms", sql, cost)
	}
	return
}

func (q *Query) Ping() error {
	return q.db.Ping()
}

func (q *Query) Exec(sql string, args ...interface{}) (result sql.Result, err error) {
	cost := utils.FuncCost(func() {
		result, err = q.db.Exec(sql, args...)
		if err != nil {
			err = errors.New(fmt.Sprintf("%s | %s", err, sql))
		}
	})
	if cost > longQueryTime {
		log.Warnf("sql: %s, cost: %dms", sql, cost)
	}
	return
}

func (q *Query) InsertMap(table string, data map[string]interface{}) (lastInsertId int64, err error) {
	fields, values := fieldsAndValues(data)
	querySql := fmt.Sprintf("INSERT INTO `%s` SET %s", table, strings.Join(fields, ","))
	var result sql.Result
	if result, err = q.db.Exec(querySql, values...); err == nil {
		lastInsertId, err = result.LastInsertId()
	}
	return
}

func (q *Query) InsertStruct(table string, stru interface{}) (lastInsertId int64, err error) {
	var data map[string]interface{}
	data, err = utils.ToMap(&stru, "orm")
	if err != nil {
		return
	}
	fields, values := fieldsAndValues(data)
	querySql := fmt.Sprintf("INSERT INTO `%s` SET %s", table, strings.Join(fields, ","))
	var result sql.Result
	if result, err = q.db.Exec(querySql, values...); err == nil {
		lastInsertId, err = result.LastInsertId()
	}
	return
}

func (q *Query) Update(table, condition string, data map[string]interface{}) (rowsAffected int64, err error) {
	if len(condition) == 0 {
		return 0, errors.New("no condition for update")
	}
	fields, values := fieldsAndValues(data)
	querySql := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", table, strings.Join(fields, ","), condition)
	var result sql.Result
	if result, err = q.db.Exec(querySql, values...); err == nil {
		rowsAffected, err = result.RowsAffected()
	}
	return
}

func (q *Query) Delete(table, condition string) (rowsAffected int64, err error) {
	if len(condition) == 0 {
		return 0, errors.New("no condition for del")
	}
	querySql := fmt.Sprintf("DELETE FROM `%s` WHERE %s", table, condition)
	var result sql.Result
	if result, err = q.db.Exec(querySql); err == nil {
		rowsAffected, err = result.RowsAffected()
	}
	return
}

func (q *Query) TruncateTable(table string) error {
	_, err := q.Query(fmt.Sprintf("TRUNCATE TABLE `%s`", table))
	return err
}

func (q *Query) DropTableIfExist(table string) error {
	_, err := q.Query(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table))
	return err
}

func (q *Query) CopyTable(tableName, newTableName string, copyContent bool) error {
	if len(newTableName) == 0 {
		return errors.New("new table name is empty")
	}
	var err error
	if _, err = q.Exec(fmt.Sprintf("CREATE TABLE `%s` LIKE `%s`", newTableName, tableName)); err == nil {
		if err = q.TruncateTable(newTableName); err == nil && copyContent {
			_, err = q.Query(fmt.Sprintf("INSERT INTO `%s` SELECT * FROM `%s`", newTableName, tableName))
		}

	}
	return err

}

func fieldsAndValues(data map[string]interface{}) (fields []string, values []interface{}) {
	for k, v := range data {
		fields = append(fields, fmt.Sprintf("`%s`=?", k))
		values = append(values, v)
	}
	return
}

func (q *Query) Replace(table string, data map[string]interface{}) (lastInsertId, rowsAffected int64, err error) {
	fields, values := fieldsAndValues(data)
	fields2 := make([]string, 0, len(fields))
	for k, _ := range data {
		fields2 = append(fields2, fmt.Sprintf("`%s`=values(`%s`)", k, k))
	}
	querySql := fmt.Sprintf("INSERT INTO `%s` SET %s ON DUPLICATE KEY UPDATE %s", table, strings.Join(fields, ","), strings.Join(fields2, ","))
	var result sql.Result
	if result, err = q.db.Exec(querySql, values...); err == nil {
		if lastInsertId, err = result.LastInsertId(); err == nil {
			rowsAffected, err = result.RowsAffected()
		}
	}
	return
}

func (q *Query) BatchInsert(table string, dataList []map[string]interface{}, ignoreDup bool) (rowsAffected int64, err error) {
	dataNum := len(dataList)
	if dataNum == 0 {
		return 0, errors.New("dataList is empty for batch replace")
	}
	one := dataList[0]
	mapLen := len(one)
	keys := make([]string, 0, mapLen)
	for k, _ := range one {
		keys = append(keys, k)
	}
	padList := make([]string, 0, dataNum)
	values := make([]interface{}, 0)
	for _, data := range dataList {
		padList = append(padList, fmt.Sprintf("(%s)", strings.TrimRight(strings.Repeat("?,", mapLen), ",")))
		for _, k := range keys {
			if v, found := data[k]; found {
				values = append(values, v)
			} else {
				return 0, errors.New("dataList is err for batch replace")
			}
		}
	}
	fields := make([]string, 0, mapLen)
	for _, k := range keys {
		fields = append(fields, fmt.Sprintf("`%s`", k))
	}
	ignore := ""
	if ignoreDup {
		ignore = "IGNORE"
	}
	querySql := fmt.Sprintf("INSERT %s INTO `%s` (%s) VALUES %s", ignore, table, strings.Join(fields, ","), strings.Join(padList, ","))
	var result sql.Result
	if result, err = q.Exec(querySql, values...); err == nil {
		rowsAffected, err = result.RowsAffected()
	}
	return
}

func (q *Query) BatchReplace(table string, dataList []map[string]interface{}, ignoreFields ...string) (lastInsertId, rowsAffected int64, err error) {
	dataNum := len(dataList)
	if dataNum == 0 {
		return 0, 0, errors.New("dataList is empty for batch replace")
	}
	one := dataList[0]
	mapLen := len(one)
	keys := make([]string, 0, mapLen)
	for k, _ := range one {
		keys = append(keys, k)
	}
	padList := make([]string, 0, dataNum)
	values := make([]interface{}, 0)
	for _, data := range dataList {
		padList = append(padList, fmt.Sprintf("(%s)", strings.TrimRight(strings.Repeat("?,", mapLen), ",")))
		for _, k := range keys {
			if v, found := data[k]; found {
				values = append(values, v)
			} else {
				return 0, 0, errors.New("dataList is err for batch replace")
			}
		}
	}

	ignoreFieldsMap := make(map[string]int, len(ignoreFields))
	if ignoreFields != nil {
		for _, f := range ignoreFields {
			ignoreFieldsMap[f] = 1
		}
	}
	fields := make([]string, 0, mapLen)
	fields2 := make([]string, 0, mapLen)
	for _, k := range keys {
		fields = append(fields, fmt.Sprintf("`%s`", k))
		if _, found := ignoreFieldsMap[k]; !found {
			fields2 = append(fields2, fmt.Sprintf("`%s`=values(`%s`)", k, k))
		}
	}
	querySql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s ON DUPLICATE KEY UPDATE %s", table, strings.Join(fields, ","), strings.Join(padList, ","), strings.Join(fields2, ","))
	var result sql.Result
	if result, err = q.Exec(querySql, values...); err == nil {
		if lastInsertId, err = result.LastInsertId(); err == nil {
			rowsAffected, err = result.RowsAffected()
		}
	}
	return
}

func (q *Query) GetDbTime() (timestamp int, err error) {
	var val string
	err = q.db.QueryRow("select unix_timestamp()").Scan(&val)
	if err == nil {
		timestamp, err = strconv.Atoi(val)
	}
	return timestamp, err
}

func rowsToVal(rows *sql.Rows, val interface{}) error {
	for rows.Next() {
		rows.Scan(val)
		break
	}
	return nil
}

func rowsToColumn(rows *sql.Rows, col interface{}) error {
	tp := reflect.TypeOf(col)
	if tp.Kind() != reflect.Ptr || tp.Elem().Kind() != reflect.Slice {
		return errors.New("list must be slice ptr")
	}
	tp = tp.Elem().Elem()
	resValue := reflect.ValueOf(col).Elem()
	for rows.Next() {
		value := reflect.New(tp)
		rows.Scan(value.Interface())
		resValue = reflect.Append(resValue, value.Elem())
	}
	reflect.ValueOf(col).Elem().Set(resValue)
	return nil
}

func rowsToMap(rows *sql.Rows) ([]map[string]string, error) {
	res := make([]map[string]string, 0)
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		m := make(map[string]string)
		values := make([]interface{}, 0, len(cols))

		for range cols {
			var v string
			values = append(values, &v)
		}
		rows.Scan(values...)
		for i, col := range cols {
			if v, ok := values[i].(*string); ok {
				m[col] = *v
			}
		}
		res = append(res, m)
	}
	return res, nil
}

func (q *Query) TableExist(table string) (res bool, err error) {
	var rows *sql.Rows
	if rows, err = q.Query(fmt.Sprintf("SHOW TABLES LIKE \"%s\"", table)); err == nil {
		for rows.Next() {
			var v string
			rows.Scan(&v)
			res = v != ""
		}
	}
	return
}
