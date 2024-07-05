package types

const (
	UnsOp    = "Unsupported operation"
	FuncProb = "Unsupported operation"
	Cant     = "Can't do this operation"
	Blank    = "Nothing to add"
	Unique   = "Not a unique value"
	LOGE     = "Problem with loging process"
	JWT      = "JWT token problem"
	AUTH     = "Not authorized"
	VALID    = "Token not valid"
	LOGGED   = "You are already loged in "
)

type Error struct {
	Error string `json:"error"`
}
