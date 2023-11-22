package dbo

const (
	SqlScriptFileExt                     = "sql"
	FileExtensionSeparator               = "."
	TableNamePrefixSeparator             = "_"
	SearchPattern_CreateTableIfNotExists = "CREATE TABLE IF NOT EXISTS %s"
	LastInsertedIdOnError                = -1
	CountOnError                         = -1
	IdOnError                            = 0
)

const (
	ErrTableNameIsNotFound = "table name is not found"
	ErrfRowsAffectedCount  = "affected rows count error: %v vs %v"
)
