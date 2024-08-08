package api

import (
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/janicaleksander/BeMotivated/auth"
	"github.com/janicaleksander/BeMotivated/components"

	"github.com/janicaleksander/BeMotivated/types"

	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Render(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {

	if r.Method == http.MethodGet {
		logged, _ := s.IsLogged(r)

		if logged {
			Render(w, r, components.SendErrorCode(1))
		}
		Render(w, r, components.RegisterForm())
		return nil
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		nickname := r.FormValue("nickname")
		password := r.FormValue("password")
		rePassword := r.FormValue("re-password")

		if password != rePassword {
			Render(w, r, components.SendErrorCode(2))
			Render(w, r, components.RegisterForm())
			return nil
		}

		acc := types.NewAccount(nickname, email, password)
		err := s.Store.CreateAccount(acc)

		if err != nil {
			fmt.Println(err)
			Render(w, r, components.SendErrorCode(3))
			Render(w, r, components.RegisterForm())
			return nil
		}

		http.Redirect(w, r, "/api/login", http.StatusSeeOther)

		return nil
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return nil
}

type LoginInformation struct {
	email    string
	password string
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if ok, _ := s.IsLogged(r); ok {
		Render(w, r, components.LoginForm("U are already loged in"))
		return nil
	}

	if r.Method == http.MethodGet {
		Render(w, r, components.LoginForm(""))
		return nil
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		linfo := LoginInformation{
			email:    email,
			password: password,
		}

		accLog := new(types.Account)
		accLog.Email = linfo.email
		accLog.Password = linfo.password
		account, err := s.Store.GetAccount(accLog.Email, accLog.Password)
		if err != nil {
			Render(w, r, components.LoginForm("Wrong email or password"))
			//	tmpl.Execute(w, FormData{Success: false, Error: "Invalid email or password", Email: email, Password: password})
			return nil
		}
		id := strconv.Itoa(account.ID)
		err = auth.CreateJWTTokenCookieUser(w, id)
		if err != nil {
			fmt.Println(err)
			//tmpl.Execute(w, FormData{Success: false, Error: "Failed to create JWT token", Email: email, Password: password})
			return nil
		}
		email, password = "", ""

		//	tmpl.Execute(w, FormData{Success: true, Error: "", Email: email, Password: password})h
		http.Redirect(w, r, "/api/dashboard", http.StatusSeeOther)

		return nil
		//return WriteToJson(w, http.StatusOK, "Login successfully")
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return nil

}

func (s *APIServer) handleTestDashboard(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}
	flag, err := s.IsLogged(r)
	if err != nil {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}
	if !flag {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.AUTH})
	}
	id, err := s.getID(r)
	if err != nil {
		return err
	}
	slice, err := s.Store.GetTask(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	Render(w, r, components.Dashboard(slice))

	return nil

}

func (s *APIServer) handleLogOut(w http.ResponseWriter, r *http.Request) error {
	/*	if r.Method != http.MethodPost {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}*/
	cookie, err := r.Cookie("jwt-token")
	if err != nil {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}
	auth.DefaultCookie(cookie)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/api/login", 302)
	return nil

}

func (s *APIServer) IsLogged(r *http.Request) (bool, error) {
	cookie, err := r.Cookie("jwt-token")
	if err != nil {
		return false, err
	}
	tokenStr := cookie.Value
	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, err
	}

	return true, nil

}

func (s *APIServer) getID(r *http.Request) (int, error) {
	cookie, err := r.Cookie("jwt-token")
	if err != nil {
		return -1, err
	}
	tokenStr := cookie.Value
	claims := &auth.Claims{}
	_, err = jwt.ParseWithClaims(tokenStr, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *APIServer) handleAddTask(w http.ResponseWriter, r *http.Request) error {
	flag, err := s.IsLogged(r)
	if err != nil {
		return WriteToJson(w, http.StatusUnauthorized, types.Error{Error: types.FuncProb})
	}
	if !flag {
		return WriteToJson(w, http.StatusUnauthorized, types.Error{Error: types.AUTH})
	}
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)

	}
	id, err := s.getID(r)
	if err != nil {
		fmt.Println(err)
		return err // error in template in the future
	}
	desc := r.PostFormValue("desc")
	dtable := strings.Split(r.PostFormValue("date"), "-")

	year := dtable[0]
	yearInt, err := strconv.Atoi(year)
	var monthInt int
	var dayInt int
	month := dtable[1]
	m := strings.Split(month, "0")
	if len(m) == 2 {
		monthInt, err = strconv.Atoi(m[1])
	} else {
		monthInt, err = strconv.Atoi(m[0])
	}

	day := dtable[2]
	d := strings.Split(day, "0")
	if len(d) == 2 {
		dayInt, err = strconv.Atoi(d[1])
	} else {
		dayInt, err = strconv.Atoi(d[0])
	}

	category := r.FormValue("cat")

	task := &types.Task{
		UserID:      id,
		Description: desc,
		Category:    category,
		CreatedAt:   time.Now(),
		Date:        time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.UTC),
	}

	err = s.Store.CreateTask(task)

	if err != nil {
		fmt.Println(err)
		return WriteToJson(w, http.StatusBadRequest, "err") // err in templ
	}
	_, err = s.Store.GetTask(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	slice, err := s.Store.GetTaskByDate(id, time.Date(2024, time.Month(7), 25, 0, 0, 0, 0, time.UTC))
	if err != nil {
		return err
	}
	for _, value := range slice {
		fmt.Println(value)
	}

	//return nil

	return Render(w, r, components.TaskInfo(*task))
	//return Render(w, r, components.Tmp())

}

func (s *APIServer) handleDeleteTask(w http.ResponseWriter, r *http.Request) error {
	tmp := r.FormValue("delete")
	item_id, err := strconv.Atoi(tmp)
	if err != nil {
		return err
	}
	err = s.Store.DeleteTask(item_id)
	if err != nil {
		return err
	}
	return nil
}

func (s *APIServer) handleTestTasks(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetTask(w http.ResponseWriter, r *http.Request) error {

	id, err := s.getID(r)
	if err != nil {
		return err
	}
	data := r.FormValue("show")

	dtable := strings.Split(data, ",")
	year := dtable[0]
	yearInt, err := strconv.Atoi(year)

	month := dtable[1]
	monthInt, err := strconv.Atoi(month)

	day := dtable[2]
	dayInt, err := strconv.Atoi(day)

	date := time.Date(yearInt, time.Month(monthInt), dayInt, 0, 0, 0, 0, time.UTC)

	sl, err := s.Store.GetTaskByDate(id, date)
	if err != nil {
		return err
	}

	Render(w, r, components.DayTaskSlice(sl))

	return nil
}

func (s *APIServer) handleDashboard(w http.ResponseWriter, r *http.Request) error {
	id, err := s.getID(r)
	if err != nil {
		return err
	}
	slice, err := s.Store.GetTaskByDate(id, time.Now())
	Render(w, r, components.DashboardProduction(slice))

	return nil

}
func (s *APIServer) TestChart(w http.ResponseWriter, r *http.Request) error {
	id, err := s.getID(r)
	if err != nil {
		fmt.Println(err)
		return err
	}
	slice := s.Store.CountCategory(id)
	fmt.Println(slice)

	return json.NewEncoder(w).Encode(slice)
}
func (s *APIServer) TestChart3(w http.ResponseWriter, r *http.Request) error {
	id, err := s.getID(r)
	if err != nil {
		fmt.Println(err)
		return err
	}
	slice, err := s.Store.CountDailyStreak(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return json.NewEncoder(w).Encode(slice)
}
func (s *APIServer) handlePomodoro(w http.ResponseWriter, r *http.Request) error {
	flag, err := s.IsLogged(r)
	if err != nil {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}
	if !flag {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.AUTH})
	}
	Render(w, r, components.PomodoroProduction())
	return nil
}
func (s *APIServer) handleSetPomodoro(w http.ResponseWriter, r *http.Request) error {

	flag, err := s.IsLogged(r)
	if err != nil {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}
	if !flag {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.AUTH})
	}
	id, err := s.getID(r)
	if err != nil {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}
	if err := s.Store.IncreasePomodoro(id); err != nil {
		fmt.Println(err)
	}

	return nil
}

func (s *APIServer) handleTask(w http.ResponseWriter, r *http.Request) error {
	id, err := s.getID(r)
	if err != nil {
		return err
	}
	slice, err := s.Store.GetTask(id)

	Render(w, r, components.SendSlice(slice))
	Render(w, r, components.TasksProduction())

	return nil

}
func (s *APIServer) handleProfile(w http.ResponseWriter, r *http.Request) error {
	flag, err := s.IsLogged(r)
	if err != nil {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}
	if !flag {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.AUTH})
	}

	Render(w, r, components.Profile())

	return nil

}
