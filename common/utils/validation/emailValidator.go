package validation

import (
	"errors"
	"strings"
)

var blackDomainsList = map[string]bool{
	"mailinator.com":    true,
	"tempmail.com":      true,
	"10minutemail.com":  true,
	"yopmail.com":       true,
	"guerrillamail.com": true,
	"trashmail.com":     true,
	"temp-mail.org":     true,
	"throwawaymail.com": true,
	"getairmail.com":    true,
	"dispostable.com":   true,
}

var whiteEmailList = map[string]bool{
	"gmail.com": true,
	"mail.ru":   true,
	"bk.ru":     true,
	"list.ru":   true,
	"inbox.ru":  true,
}

var allowedDomains = map[string]bool{ // если нужен белый список
	"mycompany.com":      true,
	"partner-company.ru": true,
}

func IsValidEmail(email string, ruDomains bool) (bool, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	at := strings.Index(email, "@")
	if at <= 0 || at >= len(email)-3 {
		return false, errors.New("incorrect format")
	}

	localPart := email[:at]
	domain := email[at+1:]

	// Базовые проверки
	if len(localPart) == 0 || len(domain) < 4 {
		return false, errors.New("local address or domain is too short")
	}
	if strings.ContainsAny(localPart, " \t\n\r") {
		return false, errors.New("spaces in the local part")
	}
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return false, errors.New("you can't start/end a period in the local part")
	}

	// 1. Черный список (disposable / временные почты)
	if blackDomainsList[domain] {
		return false, errors.New("temporary/disposable mail")
	}

	// 2. Только ру почты
	if ruDomains {
		if !allowedDomains[domain] {
			if !whiteEmailList[domain] {
				return false, errors.New("email is not available for registration")
			}
			return false, errors.New("domain name not allowed")
		}
	}

	return true, nil
}
