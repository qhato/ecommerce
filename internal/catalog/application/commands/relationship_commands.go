package commands

type CreateProductRelationshipCommand struct {
	ProductID        int64
	RelatedProductID int64
	RelationshipType string // CROSS_SELL, UP_SELL, RELATED
	Sequence         int
}

type UpdateRelationshipSequenceCommand struct {
	ID       int64
	Sequence int
}

type ActivateRelationshipCommand struct {
	ID int64
}

type DeactivateRelationshipCommand struct {
	ID int64
}

type DeleteRelationshipCommand struct {
	ID int64
}
