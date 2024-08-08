package giterror

type GitError struct {
	message string
	wrapped error
}

func NewGitError(message string, wrapped error) *GitError {
	return &GitError{message: message, wrapped: wrapped}
}

func (e *GitError) Unwrap() error {
	return e.wrapped
}

func (e *GitError) Error() string {
	return e.message
}
