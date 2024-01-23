package godotenv

type Env struct {
	Data map[string]EnvEntry
	Keys []string
}

type EnvEntry struct {
	Data    string
	Comment *string
	Quoted  bool
}
