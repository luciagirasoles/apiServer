package handlers

import (
	helpers "apiServer/helpers"
	models "apiServer/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	_ "github.com/lib/pq"
)

const (
	dbName = "apiserver"
	dbUser = "maxroach"
	dbPass = ""
	dbHost = "0.0.0.0"
	dbPort = "26257"
)

var db *sql.DB

func init() {
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	db, err = sql.Open("postgres", dbSource)

	helpers.Catch(err)

}

//AllSearches return all records
func AllSearches(w http.ResponseWriter, r *http.Request) {
	errors := []error{}
	// payload := []Domain{}
	domainList := models.DomainList{}
	stringListDomains := []string{}

	rows, err := db.Query("SELECT domain FROM domains order by ID desc;")
	helpers.Catch(err)

	defer rows.Close()

	for rows.Next() {
		data := models.Domain{}

		er := rows.Scan(&data.Domain)
		if er != nil {
			errors = append(errors, er)
			helpers.Catch(er)
		}
		// payload = append(payload, data)
		stringListDomains = append(stringListDomains, data.Domain+" info")

	}
	domainList.Items = stringListDomains
	helpers.RespondwithJSON(w, http.StatusOK, domainList)
}

//NewSearch returns all info given an domain
func NewSearch(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	// url := "truora.com/"
	data := models.Domain{}
	err := decoder.Decode(&data)

	if err != nil {
		helpers.RespondwithError(w, http.StatusNotFound, "Domain key not found")
		helpers.Catch(err)

	} else {
		if data.Domain != "" {
			query, err := db.Prepare("INSERT INTO domains (domain) VALUES ($1)")
			helpers.Catch(err)
			_, er := query.Exec(data.Domain)
			helpers.Catch(er)
			defer query.Close()
			println("Domain Saved")
		} else {
			helpers.RespondwithError(w, http.StatusNotFound, "Domain key is empty string")
			return
		}

		u, _ := url.Parse(data.Domain)
		u.Scheme = "http"
		payload := models.Information{}
		var infoErr error

		payload.Title, payload.Logo, payload.IsDown, infoErr = helpers.ElementInfo(u.String())
		if infoErr != nil {
			helpers.RespondwithError(w, http.StatusNotFound, fmt.Sprintf("%v", infoErr))
			return
		}
		payload.Servers, payload.SslGrade, payload.PreviousSslGrade, payload.ServersChanged, err = helpers.SslInfo(u.Path)
		if err == nil {
			helpers.RespondwithJSON(w, http.StatusOK, payload)
		} else {

			helpers.RespondwithError(w, http.StatusNotFound, "Error on Request")
		}

	}

}
