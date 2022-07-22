package service

type Access struct {
	Base
}

func (m *Access) GetInfo(id int) int64 {
	var accessList Access
	result := InitDB().First(&accessList, 1)
	return result.RowsAffected
}

func (m *Access) GetAllCont() int64 {
	var accessList []Access
	result := InitDB().Find(&accessList)
	return result.RowsAffected
}
