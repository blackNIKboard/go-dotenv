package godotenv

type Env struct {
	data map[string]EnvEntry
	keys []string
}

type EnvEntry struct {
	Data    string
	Comment *string
	Quoted  bool
}
