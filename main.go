package main

import(
	"database/sql"
	"encoding/json"
	"log"
    "fmt"
	"net/http"
    "strconv"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
}

type Peserta struct {
    MemberId            string  `json:"memberId"`
    LoanNo              string  `json:"loanNo"`
    // product             string  `json:"product"`
    IdentityNo          string  `json:"identityNo"`
    Name                string  `json:"name"` 
    BirthDate           string  `json:"birthDate"` 
    Gender              string  `json:"gender"` 
    Address             *string  `json:"address"` 
    // insurance           string  `json:"companyname"` 
    InsuranceStart      string  `json:"insuranceStart"` 
    InsuranceEnd        string  `json:"insuranceEnd"` 
    InsurancePeriod     int     `json:"insurancePeriod"` 
    InsuredValue        float32 `json:"insuredValue"` 
    Rate                float32 `json:"rate"` 
    Premium             float32 `json:"premium"` 
    Status              string  `json:"status"` 
    // polis               string  `json:"polis"` 
}

var db *sql.DB

func main() {
	var err error
	db,err = sql.Open("mysql","root:@tcp(127.0.0.1:3306)/bios")
    // db,err = sql.Open("mysql","root:Ad0nai_2025!@tcp(192.168.17.5:3306)/bios")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/users",getUsers).Methods("GET")
	router.HandleFunc("/user/{id}",getUser).Methods("GET")
	router.HandleFunc("/user",createUser).Methods("POST")
	router.HandleFunc("/user/{id}",updateUser).Methods("PUT")
	router.HandleFunc("/user/{id}",deleteUser).Methods("DELETE")
    router.HandleFunc("/peserta",getPeserta).Methods("GET")

    fmt.Println("starting web server at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8000",router))
}

func getUsers(w http.ResponseWriter, r *http.Request){
	var users []User
	rows, err := db.Query("SELECT id,name,email,created_at FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next(){
		var user User
		err := rows.Scan(&user.ID,&user.Name,&user.Email,&user.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users,user)
	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]
    var user User
    err := db.QueryRow("SELECT id, name, email, created_at FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            http.NotFound(w, r)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    id, err := result.LastInsertId()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    user.ID = int(id)
    user.CreatedAt = "now" // Placeholder
    json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]
    userID, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    var user User
    err = json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    user.ID = userID
    user.CreatedAt = "now" // Placeholder
    json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id := params["id"]
    userID, err := strconv.Atoi(id)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("DELETE FROM users WHERE id = ?", userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func getPeserta(w http.ResponseWriter, r *http.Request) {
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

    var pesertas []Peserta
    for rows.Next() {
        var peserta Peserta
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