package ldap

import (
	"fmt"
	"log"
	"os"
	"time"

	ldap "gopkg.in/ldap.v3"
)

const (
	ldapDateLayout      string = "2006/01/02 15:04:05"
	ldapSearchTimeLimit int    = 300 // Limit of time in second, for a ldap research
	ldapSearchSizeLimit int    = 0   // Limit of result return for a search where 0 means no limit
)

var (
	bindUsername string = os.Getenv("LDAP_BINDUSER")
	bindPassword string = os.Getenv("LDAP_BINDPASSWORD")
	url          string = os.Getenv("LDAP_URL")
	port         string = os.Getenv("LDAP_PORT")
	protocol     string = os.Getenv("LDAP_PROTOCOL")
	groupBaseDN  string = os.Getenv("LDAP_GROUPBASEDN")
)

func Users(seniorityCriteriaDate, memberOfGroup string) (seniorAccounts, newAccounts []map[string]string) {

	l, err := ldap.DialURL(fmt.Sprintf("%v://%s:%v", protocol, url, port))

	if err != nil {
		log.Println(err)
	}

	defer l.Close()

	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		log.Println(err)
	}

	searchRequest := ldap.NewSearchRequest(
		groupBaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		ldapSearchSizeLimit,
		ldapSearchTimeLimit,
		false,
		fmt.Sprintf("(&(objectClass=groupOfNames)(cn=%s))", memberOfGroup),
		[]string{"member"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Println(err)
	}

	if len(sr.Entries) != 1 {
		log.Println("Something went wrong with the group search:", memberOfGroup)
	}

	for id, member := range sr.Entries[0].Attributes[0].Values {
		// In jenkins ldap database, carLicense contains the account creation date and employeeNumber the github_id
		log.Printf("User #%v: %v\n", id+1, member)

		searchRequest = ldap.NewSearchRequest(
			member,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			ldapSearchSizeLimit,
			ldapSearchTimeLimit,
			false,
			"(&(objectClass=inetOrgPerson))",
			[]string{"cn", "mail", "carLicense", "employeeNumber"},
			nil,
		)
		sr, err = l.Search(searchRequest)
		if err != nil {
			log.Print(err)
			continue
		}

		cn := sr.Entries[0].GetAttributeValue("cn")
		mail := sr.Entries[0].GetAttributeValue("mail")
		creationDate := sr.Entries[0].GetAttributeValue("carLicense")       // Contains creation_date
		githubUsername := sr.Entries[0].GetAttributeValue("employeeNumber") // Contains github_id

		date, err := time.Parse(ldapDateLayout, creationDate)
		if err != nil && creationDate != "" {
			log.Print(err)
		}

		pivotDate, err := time.Parse(ldapDateLayout, seniorityCriteriaDate)
		if err != nil {
			log.Print(err)

		}
		// Only reject users with an account created after seniorityCriteriaDate and we assume that
		// users without creation_date were created before 2015/22/11 (https://git.io/JeGCl)

		if date.After(pivotDate) {
			newAccounts = append(newAccounts, map[string]string{"ldap_account": cn, "mail": mail, "creationDate": creationDate, "github_account": githubUsername})
			continue
		}

		seniorAccounts = append(seniorAccounts, map[string]string{"ldap_account": cn, "mail": mail, "creationDate": creationDate, "github_account": githubUsername})
	}
	return seniorAccounts, newAccounts
}
