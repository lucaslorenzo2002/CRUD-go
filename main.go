package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func connectionDB() (connection *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := "Marruecos02"
	Nombre := "CRUDgo"

	connection, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return connection
}

var templates = template.Must(template.ParseGlob("views/*"))

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)

	log.Println("Server listening...")

	http.ListenAndServe(":8080", nil)

}

type Employee struct {
	Id    int
	Name  string
	Email string
}

func Home(w http.ResponseWriter, r *http.Request) {

	establishedConnection := connectionDB()
	data, err := establishedConnection.Query("SELECT * from employees")
	if err != nil {
		panic(err.Error())
	}

	employee := Employee{}
	employeeArray := []Employee{}

	for data.Next() {
		var id int
		var name, email string
		err = data.Scan(&id, &name, &email)

		if err != nil {
			panic(err.Error())
		}

		employee.Id = id
		employee.Name = name
		employee.Email = email

		employeeArray = append(employeeArray, employee)
	}

	templates.ExecuteTemplate(w, "home", employeeArray)
}

func Create(w http.ResponseWriter, r *http.Request) {

	templates.ExecuteTemplate(w, "create", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")

		establishedConnection := connectionDB()
		insertData, err := establishedConnection.Prepare("INSERT INTO employees(name, email) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insertData.Exec(name, email)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	employeeId := r.URL.Query().Get("id")

	establishedConnection := connectionDB()
	deleteData, err := establishedConnection.Prepare("DELETE FROM employees WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	deleteData.Exec(employeeId)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	employeeId := r.URL.Query().Get("id")

	establishedConnection := connectionDB()
	editData, err := establishedConnection.Query("SELECT * FROM employees WHERE id=?", employeeId)
	if err != nil {
		panic(err.Error())
	}

	employee := Employee{}

	for editData.Next() {
		var id int
		var name, email string
		err = editData.Scan(&id, &name, &email)

		if err != nil {
			panic(err.Error())
		}

		employee.Id = id
		employee.Name = name
		employee.Email = email

	}
	templates.ExecuteTemplate(w, "edit", employee)
}

func Update(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")

		establishedConnection := connectionDB()
		updateData, err := establishedConnection.Prepare("UPDATE employees SET name=?, email=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		updateData.Exec(name, email, id)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	}
}
