package api

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/janicaleksander/BeMotivated/auth"
	"github.com/janicaleksander/BeMotivated/types"
	"net/http"
	"os"
	"strconv"
)

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

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.Cant})
	}
	accLog := new(types.Account)
	if err := json.NewDecoder(r.Body).Decode(&accLog); err != nil {

		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.FuncProb})
	}
	account, err := s.Store.GetAccount(accLog.Email, accLog.Password)
	if err != nil {

		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.LOGE})
	}
	id := strconv.Itoa(account.ID)

	err = auth.CreateJWTTokenCookieUser(w, id)
	if err != nil {
		fmt.Println(err)
		return WriteToJson(w, http.StatusUnauthorized, types.Error{Error: types.JWT})
	}

	return WriteToJson(w, http.StatusOK, "Login successfully")
	// correct email and password
	// JWT auth

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

	return WriteToJson(w, http.StatusOK, "hello on dashboard")

}

func (S *APIServer) handleLogOut(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return WriteToJson(w, http.StatusBadRequest, types.Error{Error: types.UnsOp})
	}
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
