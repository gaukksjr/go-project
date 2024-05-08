package handlers

import (
	"encoding/json"
	"fmt"
	"go/project/models"
	"go/project/repository"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	Repo *repository.Repository
}

var sessions map[string]models.Session
var SessionId string

func init() {
	sessions = make(map[string]models.Session)
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		Repo: repo,
	}
}

func (h *Handler) MenuHandler(w http.ResponseWriter, r *http.Request) {
	menuItems, err := h.Repo.GetMenu()
	if err != nil {
		http.Error(w, "Failed to retrieve menu items", http.StatusInternalServerError)
		return
	}
	log.Println("GetMenu")

	fmt.Println(SessionId)
	_, ok := sessions[SessionId]
	if !ok {
		log.Println("Invalid session")
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	jsonMenuItems, err := json.Marshal(menuItems)
	if err != nil {
		http.Error(w, "Failed to marshal menu items to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	log.Println("Menu getted successful")
	w.Write(jsonMenuItems)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println(err)
		return
	}
	log.Println(user.Password, user.Username, user.Role)

	err = h.Repo.Register(user.Username, user.Password, user.Role)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(user.Username, user.Password)

	userId, err := h.Repo.Login(user.Username, user.Password)
	if err != nil {
		http.Error(w, "invalid login or password", http.StatusBadRequest)
	}

	role, err := h.Repo.GetRoleByUserID(userId)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	SessionId = uuid.NewString()
	session := models.Session{
		UserID:   userId,
		Username: user.Username,
		Role:     role,
	}

	sessions[SessionId] = session

	w.Header().Set("Session-ID", SessionId)
	log.Println("success login")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"session_id": SessionId})
	http.Redirect(w, r, "/main", http.StatusSeeOther)

}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(SessionId)
	_, exists := sessions[SessionId]
	if !exists {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	delete(sessions, SessionId)

	log.Println("success logout")
	w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, "/home-page", http.StatusSeeOther)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}

func (h *Handler) AddMenuItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	_, ok := sessions[SessionId]
	if !ok {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	var menuItems []models.MenuItem
	err := json.NewDecoder(r.Body).Decode(&menuItems)
	if err != nil {
		http.Error(w, "Error decoding request data: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.AddMenuItems(menuItems)
	if err != nil {
		http.Error(w, "Failed to add menu items: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d menu items added successfully", len(menuItems))))
}

func (h *Handler) SaveOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	session, ok := sessions[SessionId]
	if !ok {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	var orderItems []models.OrderItem
	err := json.NewDecoder(r.Body).Decode(&orderItems)
	if err != nil {
		http.Error(w, "Error decoding request data: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID := session.UserID

	err = h.Repo.SaveOrder(userID, orderItems)
	if err != nil {
		http.Error(w, "Failed to save order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Order saved successfully"))
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home-page", http.StatusSeeOther)
}

func (h *Handler) UpdateOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	var requestData models.UpdateOrderStatus
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	session, ok := sessions[SessionId]
	if !ok {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		log.Println("Invalid session")
		return
	}

	if session.Role != "manager" {
		http.Error(w, "didn't have the permissions", http.StatusUnauthorized)
		log.Println("Unauthorized access")
		return
	}

	layout := "2006-01-02 15:04:05"

	orderTime, err := time.Parse(layout, requestData.OrderTime)
	if err != nil {
		http.Error(w, "Invalid order_time format", http.StatusBadRequest)
		log.Println("Error parsing order_time:", err)
		return
	}
	log.Println(orderTime)

	err = h.Repo.UpdateOrderStatus(orderTime)
	if err != nil {
		http.Error(w, "Failed to update order status", http.StatusInternalServerError)
		log.Println("Failed to update order status:", err)
		return
	}

	log.Println("Update order success")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Order status updated successfully"})
}

func (h *Handler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := sessions[SessionId]
	if !ok {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}
	if session.Role != "manager" {
		http.Error(w, "didn't have the permissions", http.StatusUnauthorized)
		return
	}

	orders, err := h.Repo.GetAllOrders()
	if err != nil {
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		log.Println("Failed to retrieve orders", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log.Println("get orders success")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Failed to encode orders into JSON", http.StatusInternalServerError)
		return
	}
}
func (h *Handler) GetUsernameByIDHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserId
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	username, err := h.Repo.GetUsernameByID(user.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
		log.Println("Failed to retrieve username:", err)
		return
	}
	response := struct {
		Username string `json:"username"`
	}{
		Username: username,
	}

	w.Header().Set("Content-Type", "application/json")
	log.Println("get username success")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode username into JSON", http.StatusInternalServerError)
		return
	}
}
