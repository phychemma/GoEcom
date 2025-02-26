package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"phyEcom.com/auth"
	"phyEcom.com/cart"
	"phyEcom.com/chat"
	"phyEcom.com/chat/admin"
	"phyEcom.com/file"
	"phyEcom.com/middleware"
	ordermanagement "phyEcom.com/orderManagement"
	"phyEcom.com/paystack"
	"phyEcom.com/product"
	productorder "phyEcom.com/productOrder"
	"phyEcom.com/profile"
	"phyEcom.com/review"
	"phyEcom.com/websocket"
)

func Routes(db *gorm.DB) {
	router := mux.NewRouter()
	router.Use(middleware.EnableCORS)
	manager := websocket.NewWebSocketManager()
	reviewManager := review.NewProductReviewManager()
	router.PathPrefix("/staticproductimage").Handler(http.StripPrefix("/staticproductimage", http.FileServer(http.Dir("static/images/products/"))))
	router.HandleFunc("/register", auth.Register(db)).Methods("POST")
	router.HandleFunc("/uploadPoductsInfo", file.UploadProduct2(db)).Methods("POST")
	router.HandleFunc("/login", auth.Login(db)).Methods("POST")
	router.HandleFunc("/getproductreview/{page:[0-9]+}/{limit:[0-9]+}", review.GetReviews(db)).Methods("POST")
	router.HandleFunc("/placeproductorder", productorder.PlaceOrder(db)).Methods("POST")

	router.HandleFunc("/logout", auth.Logout).Methods("GET")
	// router.HandleFunc("/fetchcategory", product.FetchCategory(db)).Methods("GET")
	router.HandleFunc("/fetchorder", productorder.FetchUserProductOrder(db)).Methods("POST")

	router.HandleFunc("/products", product.Product(db)).Methods("GET")
	router.HandleFunc("/getuserdata", profile.GetProfileDetailes(db)).Methods("GET")

	router.HandleFunc("/loginwithgoogle", auth.HandleLoginWithGoogle)
	router.HandleFunc("/callbackfromgoogle", auth.HandleCallbackFromGoogle(db))
	router.HandleFunc("/checkusername", auth.CheckIfUserExist(db)).Methods("POST")
	router.HandleFunc("/checkemail", auth.CheckIfEmailExist(db)).Methods("POST")
	router.HandleFunc("/cmfcode", auth.VerifyCode(db)).Methods("POST")
	// router.HandleFunc("/fetchsubcategory", product.FetchSubcategory(db)).Methods("POST")
	router.HandleFunc("/getproductcount", product.GetCount(db)).Methods("GET")

	router.HandleFunc("/getproduct", product.SingleProduct(db)).Methods("POST")
	router.HandleFunc("/getproduct", product.GetSingleProduct(db)).Methods("GET")

	router.HandleFunc("/getuserproductchat", chat.GetChats(db)).Methods("POST")

	router.HandleFunc("/checkifuserisauthenticated", auth.CheckIfAuthenticated).Methods("GET")

	router.HandleFunc("/orders/made/{page:[0-9]+}/{limit:[0-9]+}", ordermanagement.GetOrderMade(db)).Methods("GET")
	router.HandleFunc("/orders/received/{page:[0-9]+}/{limit:[0-9]+}", ordermanagement.GetOrderReceived(db)).Methods("GET")
	router.HandleFunc("/orders/receivedconfirm/{page:[0-9]+}/{limit:[0-9]+}", ordermanagement.GetConfirmReceivedOrder(db)).Methods("GET")

	router.HandleFunc("/orders/delivered", ordermanagement.ConfirmDeliveredOrder(db)).Methods("POST")
	router.HandleFunc("/orders/confirm", ordermanagement.ConfirmDelivered(db)).Methods("POST")
	router.HandleFunc("/getproductchat/{productID:[0-9]+}", product.ProductChat(db)).Methods("GET")
	router.HandleFunc("/getallproductchat", chat.FetchUserChattedProduct(db))
	router.HandleFunc("/getallusersinorderofchat", admin.FectchAllChattedUser(db)).Methods("GET")

	router.HandleFunc("/users", websocket.FetcUsers(db)).Methods("GET")
	router.HandleFunc("/chat/products/{page:[0-9]+}/{limit:[0-9]+}", chat.GetProducts(db)).Methods("GET")

	router.HandleFunc("/getUserProduct/{userID:[0-9]+}", admin.GetProduct(db)).Methods("GET")

	router.HandleFunc("/api/cart", cart.AddToCart(db)).Methods("POST")
	router.HandleFunc("/api/cart", cart.GetCartItems(db)).Methods("GET")
	router.HandleFunc("/api/removecart", cart.RemoveFromCart(db)).Methods("GET")
	router.HandleFunc("/api/updatecart", cart.UpdateCart(db)).Methods("GET")

	router.HandleFunc("/chatWS", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(manager, db, w, r)
	})
	router.HandleFunc("/wsreview/{productID:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		review.HandleReviewWebsocket(reviewManager, db, w, r)
	})

	router.HandleFunc("/paystack/webhook", paystack.PaystackWebhookHandler).Methods("POST")

	protected := router.PathPrefix("/admin").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.Use(middleware.RoleMiddleware("admin"))
	protected.HandleFunc("/dashboard", auth.AdminDashboard).Methods("GET")
	go manager.Run()
	go reviewManager.Run()
	log.Println("Server is running on port 999")
	log.Fatal(http.ListenAndServe(":999", router))
}
