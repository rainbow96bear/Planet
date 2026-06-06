package model

type Reaction struct {
	BaseModel
	TaskID string       `gorm:"type:uuid;not null;index:idx_reaction_unique,unique"`
	UserID string       `gorm:"type:uuid;not null;index:idx_reaction_unique,unique"`
	Type   ReactionType `gorm:"type:varchar(20);not null;index:idx_reaction_unique,unique"`
}

type ReactionType string

const (
	ReactionTypeLike  ReactionType = "like"  // ♥ 좋아요
	ReactionTypeCheer ReactionType = "cheer" // 🔥 응원
)
