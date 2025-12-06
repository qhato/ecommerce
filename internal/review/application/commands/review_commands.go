package commands

type CreateReviewCommand struct {
	ProductID     string
	CustomerID    string
	CustomerName  string
	ReviewerEmail string
	Rating        int
	Title         string
	Comment       string
	OrderID       *string
}

type UpdateReviewCommand struct {
	ID      string
	Title   string
	Comment string
	Rating  int
}

type ApproveReviewCommand struct {
	ID string
}

type RejectReviewCommand struct {
	ID string
}

type FlagReviewCommand struct {
	ID string
}

type AddResponseCommand struct {
	ID           string
	ResponseText string
}

type MarkHelpfulCommand struct {
	ID string
}

type MarkNotHelpfulCommand struct {
	ID string
}

type DeleteReviewCommand struct {
	ID string
}
