package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"BIOSGO/models"
)

func GetPeserta(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		startDate := r.URL.Query().Get("startdate")
		endDate := r.URL.Query().Get("enddate")

		query := "SELECT idpeserta,nopinjaman,noktp,nama,tgllahir,gender,alamat,tglakad,tglakhir,tenor,plafond,premirate,totalpremi,statusaktif FROM peserta WHERE del is null"
		var args []interface{}

		if startDate != "" && endDate != "" {
			query += " AND date_format(inputdate,'%Y-%m-%d') BETWEEN ? AND ?"
			args = append(args, startDate, endDate)
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var pesertas []models.Peserta
		for rows.Next() {
			var peserta models.Peserta
			err := rows.Scan(&peserta.MemberId, &peserta.LoanNo, &peserta.IdentityNo, &peserta.Name, &peserta.BirthDate, &peserta.Gender, &peserta.Address, &peserta.InsuranceStart, &peserta.InsuranceEnd, &peserta.InsurancePeriod, &peserta.InsuredValue, &peserta.Rate, &peserta.Premium, &peserta.Status)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pesertas = append(pesertas, peserta)
		}

		totalData := len(pesertas)
		var pesertaData []byte
		if totalData == 0 {
			pesertaData = []byte("[]")
		} else {
			pesertaData, err = json.Marshal(pesertas)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseData := fmt.Sprintf(`{"total": %d, "data": %s}`, totalData, pesertaData)
		w.Write([]byte(responseData))
	}
}
