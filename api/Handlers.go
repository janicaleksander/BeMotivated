package api

import (
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/janicaleksander/BeMotivated/auth"
	"github.com/janicaleksander/BeMotivated/components"
	_ "github.com/janicaleksander/BeMotivated/components"
	"github.com/janicaleksander/BeMotivated/types"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Render(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func (s *APIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}
	accReq := new(types.Account)

	if err := json.NewDecoder(r.Body).Decode(&accReq); err != nil {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.Blank})
	}
	acc := types.NewAccount(accReq.Nickname, accReq.Email, accReq.Password)

	if err := s.Store.CreateAccount(acc); err != nil {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.Cant})
	}
	return WriteToJson(w, http.StatusOK, acc)

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
		fmt.Println(email, password)
		linfo := LoginInformation{
			email:    email,
			password: password,
		}

		accLog := new(types.Account)
		accLog.Email = linfo.email
		accLog.Password = linfo.password
		account, err := s.Store.GetAccount(accLog.Email, accLog.Password)
		if err != nil {
			Render(w, r, components.LoginForm("ERROR"))
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

		//	tmpl.Execute(w, FormData{Success: true, Error: "", Email: email, Password: password})h
		http.Redirect(w, r, "/api/dashboard", http.StatusSeeOther)

		return nil
		//return WriteToJson(w, http.StatusOK, "Login successfully")
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return nil

}

func (s *APIServer) handleDashboard(w http.ResponseWriter, r *http.Request) error {
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

func (S *APIServer) handleLogOut(w http.ResponseWriter, r *http.Request) error {
	/*	if r.Method != http.MethodPost {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}*/
	cookie, err := r.Cookie("jwt-token")
	if err != nil {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}
	auth.DefaultCookie(cookie)
	http.SetCookie(w, cookie)
	fmt.Println(cookie.Expires.UTC())
	return WriteToJson(w, http.StatusOK, "Log out")

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
	id, err := s.getID(r)
	if err != nil {
		fmt.Println(err)
		return err // error in template in the future
	}
	desc := r.PostFormValue("task")
	task := &types.Task{
		UserID:      id,
		Description: desc,
		CreatedAt:   time.Now(),
	}

	err = s.Store.CreateTask(task)
	fmt.Println("id:", task.ItemID)

	if err != nil {
		fmt.Println(err)
		return WriteToJson(w, http.StatusBadRequest, "err") // err in templ
	}
	_, err = s.Store.GetTask(id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return Render(w, r, components.Task(*task))

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

func (s *APIServer) handleTest(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, components.Base())
}
func (s *APIServer) handleTestDashboard(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, components.DashboardProduction())
}
func (s *APIServer) handleTestTasks(w http.ResponseWriter, r *http.Request) error {
	return nil
}
