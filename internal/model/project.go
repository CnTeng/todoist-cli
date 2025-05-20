package model

type ProjectListArgs struct {
	All          bool `json:"all"`
	OnlyArchived bool `json:"only_archived"`
}
