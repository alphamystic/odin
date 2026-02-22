package definers

import (
  "errors"
)

var NonActiveUser = errors.New("User is not active or has not been activated.")

var UserNotLoggedIn = errors.New("Error, User isn't logged in.")

var ExpiredToken = errors.New("User has an expired token.")

var NoClaimsError = errors.New("No claims to return")

var WrongPassword = errors.New("Wrong Password")

var ErrNoMinion = errors.New("Errorno minion does not exist.")

var  Undefined = errors.New("Undefined statement or data requested/given as input.")

var ScanDoesNotExists = errors.New("Scan Does Not Exists")

var TargetDoesNotExist = errors.New("Target Does Not Exist")

var WebDataForTargetDoesNotExist = errors.New("Web Data For Target Does Not Exist")
