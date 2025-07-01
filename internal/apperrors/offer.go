package apperrors

var ErrorDeleteScoreGroupInterest = AppError{
	Err:     nil,
	Code:    500_010,
	Message: "cannot delete score group interest",
}

var InvalidStatus = AppError{Err: nil, Code: 400_1001, Message: "invalid status"}

var AssetTypeDoesNotMatch = AppError{Err: nil, Code: 400_1002, Message: "assetType does not match"}
