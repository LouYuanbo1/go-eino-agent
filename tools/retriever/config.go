package retriever

type RetrieverConfig struct {
	K               int    `mapstructure:"k"`
	IndexName       string `mapstructure:"index_name"`
	VectorFieldName string `mapstructure:"vector_field_name"`
}
