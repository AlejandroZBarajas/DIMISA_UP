package mysql

import (
	"DIMISA/src/areas/areasApp"
	"DIMISA/src/areas/areasInfra"
	"DIMISA/src/camas/camasApp"
	"DIMISA/src/camas/camasInfra"
	"DIMISA/src/cendis/cendisApp"
	"DIMISA/src/cendis/cendisInfra"
	"DIMISA/src/claves/clavesApp"
	"DIMISA/src/claves/clavesInfra"
	"DIMISA/src/colectivos/colectivosApp"
	"DIMISA/src/colectivos/colectivosInfra"
	"DIMISA/src/core/auth"
	"DIMISA/src/cpm/cpmApp"
	"DIMISA/src/cpm/cpmInfra"
	"DIMISA/src/entradaColectivo/entradaColectivoApp"
	"DIMISA/src/entradaColectivo/entradaColectivoInfra"
	"DIMISA/src/entradas/entradasApp"
	"DIMISA/src/entradas/entradasInfra"
	"DIMISA/src/inventarioODT/inventariosApp"
	"DIMISA/src/inventarioODT/inventariosInfra"
	"DIMISA/src/salidas/salidasApp"
	"DIMISA/src/salidas/salidasInfra"
	"DIMISA/src/tipos_colectivo_salida/tiposApp"
	"DIMISA/src/tipos_colectivo_salida/tiposInfra"
	"DIMISA/src/users/userApp"
	"DIMISA/src/users/userInfra"

	"database/sql"
	"log"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func RegisterRoutes(db *sql.DB) http.Handler {

	mux := http.NewServeMux()

	loginHandler := &auth.LoginHandler{DB: db}
	registerHandler := &auth.RegisterHandler{DB: db}
	mux.Handle("/login", loginHandler)
	mux.Handle("/register", registerHandler)
	log.Println("Ruta de login registrada")

	userRepo := &userInfra.UserRepository{DB: db}

	userController := userInfra.NewUserController(
		&userApp.CreateUserUseCase{Repo: userRepo},
		&userApp.UpdateUserUseCase{Repo: userRepo},
		&userApp.DeleteUserUseCase{Repo: userRepo},
		&userApp.GetAllUsersUseCase{Repo: userRepo},
		&userApp.GetUsersByRolUseCase{Repo: userRepo},
		&userApp.GetUserByIDUseCase{Repo: userRepo},
		&userApp.GetUsersByAreaUseCase{Repo: userRepo},
		&userApp.GetUsersByCendisUseCase{Repo: userRepo},
	)

	mux.HandleFunc("/users/create", userController.CreateUserHandler)          // POST
	mux.HandleFunc("/users/update", userController.UpdateUserHandler)          // PUT
	mux.HandleFunc("/users/delete", userController.DeleteUserHandler)          // DELETE
	mux.HandleFunc("/users/all", userController.GetAllUsersHandler)            // GET
	mux.HandleFunc("/users/by-rol", userController.GetUsersByRolHandler)       // POST
	mux.HandleFunc("/users/by-id", userController.GetUserByIDHandler)          // POST
	mux.HandleFunc("/users/by-area", userController.GetUsersByAreaHandler)     // POST
	mux.HandleFunc("/users/by-cendis", userController.GetUsersByCendisHandler) // POST

	log.Println("Rutas de usuarios registradas")

	camaRepo := &camasInfra.CamaRepository{DB: db}

	createCamaUC := &camasApp.CreateCama{Repo: camaRepo}
	updateCamaUC := &camasApp.UpdateCama{Repo: camaRepo}
	deleteCamaUC := &camasApp.DeleteCama{Repo: camaRepo}
	getByAreaCamaUC := &camasApp.GetCamasByArea{Repo: camaRepo}
	enableCamaUC := &camasApp.EnableCama{Repo: camaRepo}
	disableCamaUC := &camasApp.DisableCama{Repo: camaRepo}
	createCamasRangeUC := &camasApp.CreateCamasRange{Repo: camaRepo}
	getFreeCamasByAreaUC := &camasApp.GetFreeCamasByArea{Repo: camaRepo}
	setFreeCamaUC := &camasApp.SetFreeCama{Repo: camaRepo}

	camaController := camasInfra.NewCamaController(
		createCamaUC,
		updateCamaUC,
		deleteCamaUC,
		getByAreaCamaUC,
		enableCamaUC,
		disableCamaUC,
		createCamasRangeUC,
		getFreeCamasByAreaUC,
		setFreeCamaUC,
	)

	mux.HandleFunc("/camas/create", camaController.CreateCamaHandler)         // POST
	mux.HandleFunc("/camas/update", camaController.UpdateCamaHandler)         // PUT
	mux.HandleFunc("/camas/delete", camaController.DeleteCamaHandler)         // DELETE
	mux.HandleFunc("/camas/ar", camaController.GetCamasByAreaHandler)         // POST
	mux.HandleFunc("/camas/enable", camaController.EnableCamaHandler)         // PUT
	mux.HandleFunc("/camas/disable", camaController.DisableCamaHandler)       // PUT
	mux.HandleFunc("/camas/range", camaController.CreateCamasRangeHandler)    //POST
	mux.HandleFunc("/camas/frbyar", camaController.GetFreeCamasByAreaHandler) //POST
	mux.HandleFunc("/camas/setfree", camaController.SetFreeCamaHandler)

	log.Println("Rutas de camas registradas")

	areaRepo := &areasInfra.AreasRepository{DB: db}

	createAreaUC := &areasApp.CreateAreaUseCase{
		Repo:     areaRepo,
		CamaRepo: camaRepo,
	}
	updateAreaUC := &areasApp.UpdateAreaUseCase{Repo: areaRepo}
	getAllAreaUC := &areasApp.GetAllAreasUseCase{Repo: areaRepo}
	getByIDAreaUC := &areasApp.GetAreaByIDUseCase{Repo: areaRepo}
	deleteAreaUC := &areasApp.DeleteAreaUseCase{Repo: areaRepo}
	getFreeAreasUC := &areasApp.GetFreeAreasUseCase{Repo: areaRepo}
	getAreasByCendisUC := &areasApp.GetAreasByCendisUseCase{Repo: areaRepo}

	areaController := areasInfra.NewAreasController(
		createAreaUC,
		updateAreaUC,
		getAllAreaUC,
		getByIDAreaUC,
		deleteAreaUC,
		getFreeAreasUC,
		getAreasByCendisUC,
	)

	mux.HandleFunc("/areas/create", areaController.CreateAreaHandler) // POST
	mux.HandleFunc("/areas/update", areaController.UpdateAreaHandler) // PUT
	mux.HandleFunc("/areas/delete", areaController.DeleteAreaHandler) // DELETE
	mux.HandleFunc("/areas", areaController.GetAllAreasHandler)       // GET
	mux.HandleFunc("/areas/by-id", areaController.GetAreaByIDHandler) // POST
	mux.HandleFunc("/areas/free", areaController.GetFreeAreasHandler)
	mux.HandleFunc("/areas/cendis", areaController.GetAreasByCendisHandler)
	log.Println("Rutas de áreas registradas")

	// === SALIDAS ===
	salidasRepo := &salidasInfra.SalidasRepository{DB: db}
	createSlidaUC := &salidasApp.CreateSalida{Repo: salidasRepo}
	updateSalidaUC := &salidasApp.UpdateSalida{Repo: salidasRepo}
	deleteSalidaUC := &salidasApp.DeleteSalida{Repo: salidasRepo}
	getSalidasByCendisUC := &salidasApp.GetSalidasByCendis{Repo: salidasRepo}
	getSalidasPendientesUC := &salidasApp.GetSalidasPendientes{Repo: salidasRepo}
	addToSalidaUC := &salidasApp.AddToSalida{Repo: salidasRepo}
	cerrarSalidaUC := &salidasApp.CerrarSalida{Repo: salidasRepo}

	salidasController := salidasInfra.NewSalidasController(
		createSlidaUC,
		updateSalidaUC,
		deleteSalidaUC,
		getSalidasByCendisUC,
		getSalidasPendientesUC,
		addToSalidaUC,
		cerrarSalidaUC,
	)

	mux.HandleFunc("/salidas/create", salidasController.CreateSalidaHandler)
	mux.HandleFunc("/salidas/update/", salidasController.UpdateSalidaHandler)
	mux.HandleFunc("/salidas/delete", salidasController.DeleteSalidaHandler)
	mux.HandleFunc("/salidas/cendis", salidasController.GetSalidasByCendisHandler)
	mux.HandleFunc("/salidas/abiertas", salidasController.GetSalidasPendientesHandler)
	mux.HandleFunc("/salidas/add", salidasController.AddToSalidaHandler)
	mux.HandleFunc("/salidas/close", salidasController.CerrarSalidaHandler)

	log.Println(" Rutas de salidas registradas")

	// === CENDIS ===
	cendisRepo := &cendisInfra.CendisRepository{DB: db}
	createCendisUC := &cendisApp.CreateCendisUseCase{Repo: cendisRepo}
	updateCendisUC := &cendisApp.UpdateCendisUseCase{Repo: cendisRepo}
	deleteCendisUC := &cendisApp.DeleteCendisUseCase{Repo: cendisRepo}
	getAllCendisUC := &cendisApp.GetAllCendisUseCase{Repo: cendisRepo}

	cendisController := cendisInfra.NewCendisController(
		createCendisUC,
		updateCendisUC,
		getAllCendisUC,
		deleteCendisUC,
	)

	mux.HandleFunc("/cendis/create", cendisController.CreateCendisHandler) // POST
	mux.HandleFunc("/cendis/update", cendisController.UpdateCendisHandler) // PUT
	mux.HandleFunc("/cendis/delete", cendisController.DeleteCendisHandler) // DELETE
	mux.HandleFunc("/cendis/all", cendisController.GetAllCendisHandler)    // POST

	log.Println(" Rutas de cendis registradas")

	claveRepo := &clavesInfra.ClaveRepository{DB: db}
	searchMedClaveUC := &clavesApp.SearchMedClave{Repo: claveRepo}
	searchMatClaveUC := &clavesApp.SearchMatClave{Repo: claveRepo}
	searchMedInInventoryUC := &clavesApp.SearchMedInInventory{Repo: claveRepo}
	searchMatInInventoryUC := &clavesApp.SearchMatInInventory{Repo: claveRepo}
	searchAllInInventoryUC := &clavesApp.SearchAllInInventory{Repo: claveRepo}
	searchAllClavesUC := &clavesApp.SearchAllClaves{Repo: claveRepo}

	claveController := clavesInfra.NewClaveController(searchMedClaveUC, searchMatClaveUC, searchMedInInventoryUC, searchMatInInventoryUC, searchAllInInventoryUC, searchAllClavesUC)

	mux.HandleFunc("/claves/search", claveController.SearchAllClaves)                   // GET
	mux.HandleFunc("/claves/inventory/search", claveController.SearchAllInInventory)    // GET
	mux.HandleFunc("/claves/searchMed", claveController.SearchMedClave)                 // GET
	mux.HandleFunc("/claves/searchMat", claveController.SearchMatClave)                 // GET
	mux.HandleFunc("/claves/inventory/searchMed", claveController.SearchMedInInventory) // GET
	mux.HandleFunc("/claves/inventory/searchMat", claveController.SearchMatInInventory) // GET

	log.Println("Rutas de claves registradas")

	colectivosRepo := &colectivosInfra.ColectivoRepository{DB: db}

	createColectivoUC := &colectivosApp.CreateColectivo{Repo: colectivosRepo}
	getColectivosByCendisUC := &colectivosApp.GetColectivosByCendis{Repo: colectivosRepo}
	getPendingColectivosByCendisUC := &colectivosApp.GetPendingColectivosByCendis{Repo: colectivosRepo}
	getUpdatableColectivosByCendisUC := &colectivosApp.GetUpdatableColectivosByCendis{Repo: colectivosRepo}
	addToColectivoUC := &colectivosApp.AddToColectivo{Repo: colectivosRepo}
	closeColectivoUC := &colectivosApp.CloseColectivo{Repo: colectivosRepo}

	colectivosController := colectivosInfra.NewColectivosController(
		createColectivoUC,
		getColectivosByCendisUC,
		getPendingColectivosByCendisUC,
		getUpdatableColectivosByCendisUC,
		addToColectivoUC,
		closeColectivoUC,
	)

	// === ENTRADAS ===
	entradasRepo := &entradasInfra.EntradasRepository{DB: db}
	capturarEntradaUC := &entradasApp.CapturarEntradaUseCase{Repo: entradasRepo}
	capturarInventarioUC := &entradasApp.CapturarInventarioUseCase{Repo: entradasRepo}

	entradasController := entradasInfra.NewEntradaController(capturarEntradaUC, capturarInventarioUC)

	mux.HandleFunc("/entradas/capturar", entradasController.CapturarEntrada)      // POST
	mux.HandleFunc("/entradas/inventario", entradasController.CapturarInventario) // POST

	log.Println("Rutas de entradas registradas")

	// === COLECTIVOS ===
	mux.HandleFunc("/colectivos/create", colectivosController.CreateColectivoHandler)                   // POST
	mux.HandleFunc("/colectivos/by-cendis", colectivosController.GetColectivosByCendisHandler)          // POST
	mux.HandleFunc("/colectivos/pending", colectivosController.GetPendingColectivosByCendisHandler)     // POST
	mux.HandleFunc("/colectivos/editables", colectivosController.GetUpdatableColectivosByCendisHandler) //POST
	mux.HandleFunc("/colectivos/add", colectivosController.AddToColectivoHandler)
	mux.HandleFunc("/colectivos/close", colectivosController.CloseColectivoHandler) //PUT
	log.Println("Rutas de colectivos registradas")

	tiposRepo := &tiposInfra.TiposRepository{DB: db}
	getTiposUC := &tiposApp.GetTipos{Repo: tiposRepo}
	tiposController := tiposInfra.NewTiposController(getTiposUC)

	mux.HandleFunc("/tipos", tiposController.GetTiposHandler) // GET

	inventariosRepo := &inventariosInfra.InventariosRepository{DB: db}
	getInventarioByCendisIDUC := &inventariosApp.GetInventarioByCendisID{Repo: inventariosRepo}
	getAllInventariosUC := &inventariosApp.GetAllInventarios{Repo: inventariosRepo}
	inventariosController := inventariosInfra.NewInventariosController(getInventarioByCendisIDUC, getAllInventariosUC)

	mux.HandleFunc("/inventarios/cendis", inventariosController.GetInventarioByCendisID)
	mux.HandleFunc("/inventarios/all", inventariosController.GetAllInventarios)

	cpmRepo := &cpmInfra.CpmRepository{DB: db}
	getCpmUC := &cpmApp.GetCpm{Repo: cpmRepo}
	cpmController := cpmInfra.NewCpmController(getCpmUC)

	mux.HandleFunc("/cpm", cpmController.GetCpm) // GET

	log.Println("Rutas de CPM registradas")
	//handlerWithCors := corsMiddleware(mux)

	// === ENTRADAS COLECTIVO ===
	entradaColectivoRepo := &entradaColectivoInfra.EntradaColectivoRepository{DB: db}

	getReporteMensualUC := &entradaColectivoApp.GetReporteMensual{Repo: entradaColectivoRepo}
	getDeficitCronicoUC := &entradaColectivoApp.GetDeficitCronico{Repo: entradaColectivoRepo}
	getComparativoCendisUC := &entradaColectivoApp.GetComparativoCendis{Repo: entradaColectivoRepo}

	entradaColectivoController := entradaColectivoInfra.NewEntradaColectivoController(
		getReporteMensualUC,
		getDeficitCronicoUC,
		getComparativoCendisUC,
	)

	mux.HandleFunc("/entradas-colectivo/reporte-mensual", entradaColectivoController.GetReporteMensual) // POST
	mux.HandleFunc("/entradas-colectivo/deficit-cronico", entradaColectivoController.GetDeficitCronico) // POST
	mux.HandleFunc("/entradas-colectivo/comparativo", entradaColectivoController.GetComparativoCendis)  // POST

	log.Println("Rutas de entradas colectivo registradas")

	//log.Fatal(http.ListenAndServe(":8080", handlerWithCors))
	log.Println("Servidor escuchando en :8080")
	return corsMiddleware(mux)
}
