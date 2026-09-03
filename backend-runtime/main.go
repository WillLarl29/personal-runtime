package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "backend-runtime/docs"
)

// pgUniqueViolation es el código de error de Postgres para violaciones de
// restricciones UNIQUE (aquí, UNIQUE(actividad_id, fecha) en checks_diarios).
const pgUniqueViolation = "23505"

type server struct {
	db *pgxpool.Pool
}

// actividadSelect trae cada actividad ya con el nombre de su categoría
// (LEFT JOIN porque categoria_id es opcional) y si ya se registró un check
// hoy (CURRENT_DATE, en la zona horaria del servidor de base de datos), para
// que el frontend no tenga que adivinar el estado del día.
const actividadSelect = `
	SELECT a.id, a.titulo, a.descripcion, a.categoria_id, c.nombre AS categoria_nombre,
	       a.prioridad, a.activa,
	       EXISTS (
	         SELECT 1 FROM checks_diarios ch
	         WHERE ch.actividad_id = a.id AND ch.fecha = CURRENT_DATE
	       ) AS check_hoy,
	       a.creado_en, a.actualizado_en
	FROM actividades a
	LEFT JOIN categorias c ON c.id = a.categoria_id
`

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
	mux.HandleFunc("GET /categorias", s.handleListarCategorias)
	mux.HandleFunc("POST /categorias", s.handleCrearCategoria)
	mux.HandleFunc("PUT /categorias/{id}", s.handleActualizarCategoria)
	mux.HandleFunc("DELETE /categorias/{id}", s.handleEliminarCategoria)
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
	rows, err := s.db.Query(r.Context(), actividadSelect+" WHERE a.activa = TRUE ORDER BY a.id")
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

	var nuevoID int32
	err := s.db.QueryRow(r.Context(),
		`INSERT INTO actividades (titulo, descripcion, categoria_id, prioridad)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		input.Titulo, input.Descripcion, input.CategoriaID, input.Prioridad,
	).Scan(&nuevoID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	rows, err := s.db.Query(r.Context(), actividadSelect+" WHERE a.id = $1", nuevoID)
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
// @Description  Solo se permite un check por actividad y por día (fecha del
// @Description  servidor); si ya existe uno, devuelve 409.
// @Tags         checks
// @Accept       json
// @Produce      json
// @Param        check  body      CrearCheckInput  true  "Datos del check"
// @Success      201  {object}  CheckDiario
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string  "Ya existe un check para esta actividad hoy"
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
		writeCheckError(w, err)
		return
	}
	defer rows.Close()

	check, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[CheckDiario])
	if err != nil {
		writeCheckError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, check)
}

// writeCheckError distingue la violación de la restricción "un check por
// actividad y día" (UNIQUE(actividad_id, fecha)) del resto de errores de
// base de datos, para que el frontend reciba un 409 con mensaje claro en vez
// de un 500 genérico cuando el usuario ya marcó la actividad hoy.
func writeCheckError(w http.ResponseWriter, err error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		writeError(w, http.StatusConflict, errors.New("esta actividad ya tiene un check registrado hoy"))
		return
	}
	writeError(w, http.StatusInternalServerError, err)
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

// handleListarCategorias godoc
// @Summary      Listar categorías
// @Tags         categorias
// @Produce      json
// @Success      200  {array}   Categoria
// @Failure      500  {object}  map[string]string
// @Router       /categorias [get]
func (s *server) handleListarCategorias(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(r.Context(), "SELECT id, nombre FROM categorias ORDER BY nombre")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	categorias, err := pgx.CollectRows(rows, pgx.RowToStructByName[Categoria])
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, categorias)
}

// handleCrearCategoria godoc
// @Summary      Crear categoría
// @Tags         categorias
// @Accept       json
// @Produce      json
// @Param        categoria  body      CrearCategoriaInput  true  "Nombre de la categoría"
// @Success      201  {object}  Categoria
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categorias [post]
func (s *server) handleCrearCategoria(w http.ResponseWriter, r *http.Request) {
	var input CrearCategoriaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	rows, err := s.db.Query(r.Context(),
		"INSERT INTO categorias (nombre) VALUES ($1) RETURNING id, nombre",
		input.Nombre,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	categoria, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Categoria])
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, categoria)
}

// handleActualizarCategoria godoc
// @Summary      Renombrar categoría
// @Tags         categorias
// @Accept       json
// @Produce      json
// @Param        id         path      int                       true  "ID de la categoría"
// @Param        categoria  body      ActualizarCategoriaInput  true  "Nuevo nombre"
// @Success      200  {object}  Categoria
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categorias/{id} [put]
func (s *server) handleActualizarCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var input ActualizarCategoriaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	rows, err := s.db.Query(r.Context(),
		"UPDATE categorias SET nombre = $1 WHERE id = $2 RETURNING id, nombre",
		input.Nombre, id,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	categoria, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Categoria])
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	writeJSON(w, http.StatusOK, categoria)
}

// handleEliminarCategoria godoc
// @Summary      Eliminar categoría
// @Tags         categorias
// @Param        id  path  int  true  "ID de la categoría"
// @Success      204  "Sin contenido"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categorias/{id} [delete]
func (s *server) handleEliminarCategoria(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if _, err := s.db.Exec(r.Context(), "DELETE FROM categorias WHERE id = $1", id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
