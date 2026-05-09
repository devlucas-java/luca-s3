package errors

func ErrInvalidCredentials(err error) *AppError {
	return New(
		"INVALID_CREDENTIALS",
		"Invalid email/username or password",
		401,
		err,
	)
}

func ErrInvalidRole(message string, err error) *AppError {
	return New(
		"INVALID_ROLE",
		message,
		400,
		err,
	)
}

func ErrInvalidUUID(err error) *AppError {
	return New(
		"INVALID_UUID",
		"Invalid UUID provided",
		400,
		err,
	)
}

func ErrNotFound(resource string, err error) *AppError {
	return New(
		resource+"_NOT_FOUND",
		resource+" not found",
		404,
		err,
	)
}
func ErrConflict(resource string, err error) *AppError {
	return New(
		resource+"_CONFLICT",
		resource+" already exists",
		409,
		err,
	)
}
func ErrUnauthorized(err error) *AppError {
	return New(
		"UNAUTHORIZED",
		"Unauthorized",
		401,
		err,
	)
}
func ErrForbidden(err error) *AppError {
	return New(
		"FORBIDDEN",
		"Forbidden",
		403,
		err,
	)
}
func ErrBadRequest(message string, err error) *AppError {
	return New(
		"BAD_REQUEST",
		message,
		400,
		err,
	)
}
func ErrUnprocessable(message string, err error) *AppError {
	return New(
		"UNPROCESSABLE_ENTITY",
		message,
		422,
		err,
	)
}
func ErrDatabase(message string, err error) *AppError {
	return New(
		"DATABASE_ERROR",
		message,
		500,
		err,
	)
}
func ErrInternal(message string, err error) *AppError {
	return New(
		"INTERNAL_SERVER_ERROR",
		message,
		500,
		err,
	)
}
