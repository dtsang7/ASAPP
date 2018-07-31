package models

func (dao *DAO) CheckDB() (int, error) {
	var res int
	err := dao.db.QueryRow("SELECT 1").Scan(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}
