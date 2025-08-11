package errors

import "fmt"

var (
	ErrUnknown           = fmt.Errorf("unknown error occured, read logs to learn more")
	ErrNotFound          = fmt.Errorf("entity was not found in database")
	ErrUnableToQuery     = fmt.Errorf("could not get the queried entity")
	ErrUnableToSave      = fmt.Errorf("entity could not be saved into database")
	ErrUnableToUpdate    = fmt.Errorf("entity could not be updated")
	ErrUnableToDelete    = fmt.Errorf("entity could not be removed")
	ErrNotReceivedInputs = fmt.Errorf("didn't receive enought input parameters")
)
