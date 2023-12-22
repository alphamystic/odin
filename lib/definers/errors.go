package definers

import (
  "errors"
)

var NonActiveUser = errors.New("User is not active or has not been activated.")

var UserNotLoggedIn = errors.New("Error, User isn't logged in.")

var ExpiredToken = errors.New("User has an expired token.")
