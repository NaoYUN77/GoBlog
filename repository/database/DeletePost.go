package database



func DeletePost(id int64) error {

	sqlstr := "delete from tb_post where id=?"

	_ , err :=  DB.Exec(sqlstr, id)
	return err
}
