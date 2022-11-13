package godotenv

type Env map[string]EnvEntry

type EnvEntry struct {
	Data    string
	Comment *string
	Quoted  bool
}
