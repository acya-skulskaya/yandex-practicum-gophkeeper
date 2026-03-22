package secret

import "errors"

var ErrNameAndTypeAlreadyExist = errors.New("a secret with this name and type already exists")
var ErrSecretNotFound = errors.New("could not find secret")
var ErrSecretVersionNotFound = errors.New("could not find secret version")
var ErrSecretHasOnlyOneOrLessVersions = errors.New("secret has only one or less versions")
