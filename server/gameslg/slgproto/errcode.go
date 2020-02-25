package slgproto

const (
	CodeSlgSuccess = 0
	CodeDbError    = iota + 10000
	CodeNotAuth
	CodeRoleExit
	CodeRoleNotFound
	CodeBuildingUpError
	CodeGeneralError
	CodeCityError
	CodeNotLocalCity //不是本国领土
	CodeAttackLocalCity
)

