package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "backend-runtime/docs"
)

type server struct {
	db *pgxpool.Pool
}

// @title           API de Actividades
// @version         1.0
// @description     API para gestionar actividades y sus checks diarios.
// @host            localhost:3000
// @BasePath        /
func main() {
	ctx := context.Background()

	pool, err := NewPool(ctx)
	if err != nil {
		log.Fatalf("error inicializando la base de datos: %v", err)
	}
	defer pool.Close()

	s := &server{db: pool}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handleRoot)
	mux.HandleFunc("GET /actividades", s.handleListarActividades)
	mux.HandleFunc("POST /actividades", s.handleCrearActividad)
	mux.HandleFunc("POST /checks", s.handleCrearCheck)
	mux.HandleFunc("GET /resumen", s.handleResumen)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("API en http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

// withCORS permite que el frontend (Vite, en otro puerto) llame a esta API.
// El origen permitido se lee de CORS_ORIGIN; por defecto es el puerto de
// desarrollo de Vite.
func withCORS(next http.Handler) http.Handler {
	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "http://localhost:5173"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API funcionando ✅"))
}

// handleListarActividades godoc
// @Summary      Listar actividades
// @Description  Devuelve las actividades activas, ordenadas por id
// @Tags         actividades
// @Produce      json
// @Success      200  {array}   Actividad
// @Failure      500  {object}  map[string]string
// @Router       /actividades [get]
func (s *server) handleListarActividades(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(r.Context(),
		"SELECT id, titulo, descripcion, categoria, prioridad, activa, creado_en, actualizado_en FROM actividades WHERE activa = TRUE ORDER BY id",
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	actividades, err := pgx.CollectRows(rows, pgx.RowToStructByName[Actividad])
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, actividades)
}

// handleCrearActividad godoc
// @Summary      Crear actividad
// @Tags         actividades
// @Accept       json
// @Produce      json
// @Param        actividad  body      CrearActividadInput  true  "Datos de la actividad"
// @Success      201  {object}  Actividad
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /actividades [post]
func (s *server) handleCrearActividad(w http.ResponseWriter, r *http.Request) {
	var input CrearActividadInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	rows, err := s.db.Query(r.Context(),
		`INSERT INTO actividades (titulo, descripcion, categoria, prioridad)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, titulo, descripcion, categoria, prioridad, activa, creado_en, actualizado_en`,
		input.Titulo, input.Descripcion, input.Categoria, input.Prioridad,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	actividad, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Actividad])
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, actividad)
}

// handleCrearCheck godoc
// @Summary      Dar check a una actividad hoy
// @Tags         checks
// @Accept       json
// @Produce      json
// @Param        check  body      CrearCheckInput  true  "Datos del check"
// @Success      201  {object}  CheckDiario
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /checks [post]
func (s *server) handleCrearCheck(w http.ResponseWriter, r *http.Request) {
	var input CrearCheckInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	rows, err := s.db.Query(r.Context(),
		`INSERT INTO checks_diarios (actividad_id, nota)
		 VALUES ($1, $2)
		 RETURNING id, actividad_id, fecha, nota, creado_en`,
		input.ActividadID, input.Nota,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	check, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[CheckDiario])
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, check)
}

// handleResumen godoc
// @Summary      Resumen de cumplimiento por actividad
// @Description  Usa la vista resumen_actividades
// @Tags         resumen
// @Produce      json
// @Success      200  {array}   ResumenActividad
// @Failure      500  {object}  map[string]string
// @Router       /resumen [get]
func (s *server) handleResumen(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(r.Context(), "SELECT * FROM resumen_actividades")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	resumen, err := pgx.CollectRows(rows, pgx.RowToStructByName[ResumenActividad])
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, resumen)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
