package apperrors

var ErrInvalidSymbolScore = New(nil, WithCode(400_0010), WithMessage("invalid symbol score"))
