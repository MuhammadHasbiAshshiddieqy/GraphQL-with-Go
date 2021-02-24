# GraphQL-with-Go #

# Cara Mengubah schema #

1. Ubah schema di file graph/schema.graphqls
2. run fungsin ini : rm graph/schema.resolvers.go
3. run fungsi ini di terminal go run github.com/99designs/gqlgen generate
4. buka graph/model/models_gen.go dan ubah struct `CustomerAddress` pada variabel `Address`, dari `json:"address"` menjadi `gorm:"column:address_1" json:"address"`
5. git checkout graph/schema.resolvers.go
6. Modifikasi file graph/schema.resolvers.go sesuai yang diperlukan
