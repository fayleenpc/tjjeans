package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fayleenpc/ecom/config"
	"github.com/fayleenpc/ecom/service/cart"
	"github.com/fayleenpc/ecom/service/gateway/payment"
	"github.com/fayleenpc/ecom/service/landing"
	"github.com/fayleenpc/ecom/service/order"
	"github.com/fayleenpc/ecom/service/products"
	"github.com/fayleenpc/ecom/service/user"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

// test /product post or /product get
// {
// 	"name": "new product12612512",
// 	"description": "new product arrive61251261",
// 	"image" : "new image612612512",
// 	"price": 1000000,
// 	"quantity": 300
// }

// {
// 	"id": 3,
// 	"name": "new jeans",
// 	"description": "new jeans product",
// 	"image": "image new jeans",
// 	"quantity": 100,
// 	"price": 80000
// }

// test /login
// {
//     "email": "me@me.com",
//     "password": "asd"
// }

// test /register
// {
//     "email": "me@me.com",
//     "password": "asd",
//     "firstName": "tiago",
//     "lastName": "user"
// }

// test /cart/checkout
// {
// 	"items": [
// 	  {
// 		"productID": 1,
// 		"quantity": 2
// 	  },
// 	  {
// 		"productID": 2,
// 		"quantity": 3
// 	  }
// 	]
// }

// test /payment/invoices
// {
// 	"payment" : {
// 	  "payment_type" : "alfamart"
// 	},
// 	"customer": {
// 	  "name" : "john",
// 	  "email" : "foo@bar.com",
// 	  "phone_number" : ""
// 	},
// 	"items" : [{
// 	  "name" : "support podcast",
// 	  "category" : "podcast",
// 	  "merchant": "imregi.com",
// 	  "description": "donasi podcast imre",
// 	  "qty": 1,
// 	  "price": 10000,
// 	  "currency": "IDR"
// 	},
// 	{
// 	  "name" : "gk1020h12",
// 	  "category" : "podcast gk1020h12",
// 	  "merchant": "imregi.com gk1020h12",
// 	  "description": "donasi gk1020h12",
// 	  "qty": 1,
// 	  "price": 50000,
// 	  "currency": "IDR"
// 	}]
//   }

func (s *APIServer) Run() error {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:8080/", "http://localhost:8080/api/v1/login", "http://localhost:8080/api/v1/register", "http://localhost:8080/api/v1/payment/*", "http://localhost:8080/api/v1/cart/checkout"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := products.NewStore(s.db)
	productHandler := products.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	// payment gateway
	srv := payment.NewHandler(subrouter, payment.NewServer())
	srv.RegisterRoutes()
	// landing
	landing := landing.NewHandler()
	landing.RegisterRoutes(router)

	// serve files in static folder
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.PathPrefix("/static/images/").Handler(http.StripPrefix("/static/images/", http.FileServer(http.Dir("static/images"))))

	// serve files in static-admin folder
	// router.PathPrefix("/static-admin/").Handler(http.StripPrefix("/static-admin/", http.FileServer(http.Dir("static-admin"))))
	// router.PathPrefix("/static-admin/imgs/").Handler(http.StripPrefix("/static-admin/imgs/", http.FileServer(http.Dir("static-admin/imgs"))))
	// router.PathPrefix("/static-admin/css/").Handler(http.StripPrefix("/static-admin/css/", http.FileServer(http.Dir("static-admin/css"))))
	// router.PathPrefix("/static-admin/js/").Handler(http.StripPrefix("/static-admin/js/", http.FileServer(http.Dir("static-admin/js"))))
	log.Println("Listening on ", s.addr)
	log.Println("ENVS : ")
	log.Println(config.Envs)

	return http.ListenAndServe(s.addr, handlers.CORS(originsOk, headersOk, methodsOk)(router))
}
