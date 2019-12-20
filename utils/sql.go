package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"errors"
)


type SqlCliet struct {
	DB *sql.DB
}

func InitSql()(*SqlCliet,error) {
	var err error
	var db *sql.DB
	db,err = sql.Open("mysql",DbUser+":"+DbPwd+"@tcp("+DbAddr+":"+DbPort+")/"+DbName+"?charset=utf8&parseTime=true")
	if err != nil {
		return nil,err
	}
	err = db.Ping()
	if err != nil {
		return nil,err
	}
	return &SqlCliet{
		DB:db,
	},nil
}

func (s *SqlCliet)InsertFunc(funcId string)error{
	stmt, err := s.DB.Prepare("INSERT INTO func(funcid) VALUES (?) ")
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(funcId)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlCliet)DeleteFuncById(funcId string)(int64,error){
	stmt, err := s.DB.Prepare(`DELETE FROM func WHERE funcid=?`)
	defer stmt.Close()
	if err != nil {
		return 0,err
	}
	res, err := stmt.Exec(funcId)
	if err != nil {
		return 0,err
	}

	num, err := res.RowsAffected()
	if err != nil {
		return 0,err
	}
	return num,nil
}


func (s *SqlCliet)GetFunc(funcId string)(string,error){
	stmt,err := s.DB.Prepare("select funcid from func where funcid = ?")
	defer stmt.Close()
	if err != nil {
		return "", err
	}
	row := stmt.QueryRow(funcId)

	if row != nil {
		var funid string
		err = row.Scan(&funid)
		if err != nil {
			if sql.ErrNoRows == err {
				return "",nil
			}
			return "" , err
		}
		return funid,nil
	}else{
		return "" , errors.New("the func not exist")
	}
}

func (s *SqlCliet)InsertFuncList(funcId string,ip string)error{
	stmt, err := s.DB.Prepare("INSERT INTO funclist(funcid,ip) VALUES (?,?) ")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(funcId,ip)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlCliet)DeleteFuncList(funcId,ip string)(int64,error){
	stmt, err := s.DB.Prepare(`DELETE FROM funclist WHERE funcid=? and ip =?`)
	defer stmt.Close()
	if err != nil {
		return 0,err
	}
	res, err := stmt.Exec(funcId,ip)
	if err != nil {
		return 0,err
	}

	num, err := res.RowsAffected()
	if err != nil {
		return 0,err
	}
	return num,nil
}

func (s *SqlCliet)GetFuncListByIP(funcId,ip string)(string,error){
	stmt,err := s.DB.Prepare("select ip from funclist where funcid = ? and ip=?")
	defer stmt.Close()
	if err != nil {
		return "", err
	}
	row := stmt.QueryRow(funcId,ip)
	if row != nil {
		var ip string
		err = row.Scan(&ip)
		if err != nil {
			if sql.ErrNoRows == err {
				return "",nil
			}
			return "" , err
		}
		return ip,nil
	}else{
		return "" , errors.New("the funcid in funclist table not exist")
	}
}

func (s *SqlCliet)GetFuncList(funcId string)(string,error){
	stmt,err := s.DB.Prepare("select ip from funclist where funcid = ?")
	defer stmt.Close()
	if err != nil {
		return "", err
	}
	row := stmt.QueryRow(funcId)
	if row != nil {
		var ip string
		err = row.Scan(&ip)
		if err != nil {
			if sql.ErrNoRows == err {
				return "",nil
			}
			return "" , err
		}
		return ip,nil
	}else{
		return "" , errors.New("the funcid in funclist table not exist")
	}
}

func (s *SqlCliet)GetAllFunc()(map[string]string,error){
	stmt,err := s.DB.Prepare("select funcid, ip from funclist")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	rows,err := stmt.Query()
	if err != nil {
		return nil,err
	}
	var funcipList = map[string]string{}
	for rows.Next(){
		var funcid,ip string
		err = rows.Scan(&funcid,&ip)
		if err != nil {
			continue
		}
		funcipList[funcid]=ip
	}
	return funcipList,nil
}

func (s *SqlCliet)InsertSvc(svcId string)error{
	stmt, err := s.DB.Prepare("INSERT INTO services(svcid) VALUES (?) ")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(svcId)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlCliet)GetSvc(svcId string)(string,error){
	stmt,err := s.DB.Prepare("select svcid from services where svcid = ?")
	defer stmt.Close()
	if err != nil {
		return "", err
	}
	row := stmt.QueryRow(svcId)
	if row != nil {
		var id string
		err = row.Scan(&id)
		if err != nil {
			if sql.ErrNoRows == err {
				return "",nil
			}
			return "",err
		}
		return id,nil
	}else{
		return "" , errors.New("the funcid in funclist table not exist")
	}
}

func (s *SqlCliet)GetAllSvc()([]string,error){
	stmt,err := s.DB.Prepare("select svcid from services")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	rows,err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var list = []string{}
	for rows.Next(){
		var svcno string
		err = rows.Scan(&svcno)
		if err != nil {
			continue
		}
		list = append(list,svcno)
	}
	return list,nil
}

func (s *SqlCliet)InsertSubScribe(svcId,Ip string)error{
	stmt, err := s.DB.Prepare("INSERT INTO subscribes(svcid,ip) VALUES (?,?) ")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(svcId,Ip)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlCliet)GetSubScribe(svcId string)([]string,error){
	stmt,err := s.DB.Prepare("select ip from subscribes where svcid = ?")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	rows,err := stmt.Query(svcId)
	if err != nil {
		return nil,err
	}

	listIps := make([]string,0)
	for rows.Next(){
		var ip string
		err = rows.Scan(&ip)
		if err != nil {
			continue
		}
		listIps =  append(listIps,ip)
	}
	return listIps,nil
}

func (s *SqlCliet)GetSubScribeByIP(svcId,ip string)(string,error){
	stmt,err := s.DB.Prepare("select ip from subscribes where svcid = ? and ip = ?")
	defer stmt.Close()
	if err != nil {
		return "", err
	}
	row := stmt.QueryRow(svcId,ip)
	if row != nil {
		var ip string
		err = row.Scan(&ip)
		if err != nil {
			if sql.ErrNoRows == err {
				return "",nil
			}
			return "",err
		}
		return ip,nil
	}else{
		return "" , errors.New("the svcId and ip in subscribe table not exist")
	}
}

func (s *SqlCliet)GetAllSubScribe()(map[string][]string,error){
	svcList,err := s.GetAllSvc()
	if err != nil {
		return nil,err
	}
	var subscribeDetail = map[string][]string{}
	for _,svcid := range svcList{
		listIP , err := s.GetSubScribe(svcid)
		if err != nil {
			continue
		}
		subscribeDetail[svcid] = listIP
	}
	return subscribeDetail,nil
}
// CREATE DATABASE IF NOT EXISTS midmsg DEFAULT CHARSET utf8 COLLATE utf8_general_ci;
// CREATE TABLE func( funcid VARCHAR(50) NOT NULL, PRIMARY KEY (funcid ) )ENGINE=InnoDB DEFAULT CHARSET=utf8;
// CREATE TABLE funclist( funcid VARCHAR(50) NOT NULL, ip VARCHAR(50) NOT NULL, PRIMARY KEY (funcid ) )ENGINE=InnoDB DEFAULT CHARSET=utf8;
// CREATE TABLE services( svcid VARCHAR(50) NOT NULL, PRIMARY KEY (svcid ) )ENGINE=InnoDB DEFAULT CHARSET=utf8;
// CREATE TABLE subscribes( svcid VARCHAR(50) NOT NULL,ip VARCHAR(50) NOT NULL )ENGINE=InnoDB DEFAULT CHARSET=utf8;