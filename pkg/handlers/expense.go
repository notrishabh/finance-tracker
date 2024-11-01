package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/notrishabh/finance-tracker/pkg/middleware"
	"github.com/notrishabh/finance-tracker/pkg/models"
	"github.com/notrishabh/finance-tracker/pkg/services"
)

func GetExpenseHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "User not found, please login again.", http.StatusUnauthorized)
		return
	}
	expenses, err := services.GetExpense(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(expenses)
}

func GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "User not found, please login again.", http.StatusUnauthorized)
		return
	}
	categories, err := services.GetCategory(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "User not found, please login again.", http.StatusUnauthorized)
		return
	}

	var category models.Category
	json.NewDecoder(r.Body).Decode(&category)
	err := services.CreateCategory(&category, user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(category)
}

func CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "User not found, please login again.", http.StatusUnauthorized)
		return
	}
	var rawBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&rawBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	category, ok := rawBody["category"].(string)
	if !ok {
		http.Error(w, "category is required", http.StatusBadRequest)
		return
	}

	var expense models.Expense
	expenseData, err := json.Marshal(rawBody)
	if err != nil {
		http.Error(w, "Error processing expense data", http.StatusInternalServerError)
		return
	}

	json.Unmarshal(expenseData, &expense)
	error := services.CreateExpense(&expense, category, user.Username)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(expense)
}

func UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var rawBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&rawBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var category string
	if rawBody["category"] != nil {
		category = rawBody["category"].(string)
	}

	var expense models.Expense
	expenseData, err := json.Marshal(rawBody)
	if err != nil {
		http.Error(w, "Error processing expense data", http.StatusInternalServerError)
		return
	}

	json.Unmarshal(expenseData, &expense)
	err = services.UpdateExpense(params["id"], &expense, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Expense updated successfully"))
}

func DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := services.DeleteExpense(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Expense deleted successfully"))
}

func GetExpenseByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	expense, err := services.GetExpenseById(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(expense)
}
