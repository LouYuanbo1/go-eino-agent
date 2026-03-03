package params

type PythonParams struct {
	Code string `json:"code" jsonschema:"description=用于指定要执行的 Python 代码"`
}
