package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type pgConnection struct {
	Endpoint string
	Port     string
	User     string
	Password string
	Database string
	Status   string
}

func (p pgConnection) withProperty(name string, value string) pgConnection {
	switch strings.ToLower(name) {
	case "endpoint":
		p.Endpoint = value
	case "port":
		p.Port = value
	case "user":
		p.User = value
	case "password":
		p.Password = value
	case "database":
		p.Database = value
	default:
		log.Printf("Unknown property : %s", name)
	}
	return p
}

func (p pgConnection) build() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", p.User, p.Password, p.Endpoint, p.Port, p.Database)
}

func (connection *pgConnection) tryConnect() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	conn, err := pgx.Connect(ctx, connection.build())
	if err != nil {
		connection.Status = fmt.Sprintf("unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())
	connection.Status = "connected to database"
}

type envModel struct {
	EnvVars       []string
	PgConnections []pgConnection
}

func main() {
	// Définir le template HTML
	tmpl := template.Must(template.New("env").Parse(`
		<html>
		<body>
			<h1>Env vars :</h1>
			<ul>
				{{range .EnvVars}}
					<li>{{.}}</li>
				{{end}}
			</ul>
			<h1>PostgreSQL Connections :</h1>
			<ul>
				{{range .PgConnections}}
					<li>
						{{.Endpoint}}:{{.Port}} - {{.User}} - {{.Database}} : {{.Status}}
					</li>
				{{end}}
			</ul>
		</body>
		</html>
	`))

	http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		// All vars
		envVars := os.Environ()

		mapConns := make(map[string]pgConnection)

		for _, envVar := range envVars {
			// key=value
			parts := strings.SplitN(envVar, "=", 2)
			if len(parts) != 2 {
				log.Printf("malformed env var : %s", envVar)
				continue
			}
			varName := parts[0]
			varValue := parts[1]

			regex, err := regexp.Compile(`^((API|AUTH)_(\d))_(.*)$`)
			if err != nil {
				log.Printf("regex error: %v", err)
				continue
			}

			if regex.MatchString(varName) {
				matches := regex.FindStringSubmatch(varName)
				if len(matches) >= 5 {
					idx := matches[2]
					property := matches[4]

					if _, ok := mapConns[idx]; !ok {
						mapConns[idx] = pgConnection{Status: "not connected"}
					}

					mapConns[idx] = mapConns[idx].withProperty(property, varValue)
					log.Printf("idx: %s, property: %s, value: %s", idx, property, varValue)
				}
			}
		}

		pgConns := make([]pgConnection, 0, len(mapConns))
		for _, conn := range mapConns {
			conn.tryConnect()
			pgConns = append(pgConns, conn)
		}

		model := envModel{envVars, pgConns}

		// Rendu du template avec les données
		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, model); err != nil {
			log.Printf("template render error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	})

	log.Println("server start on http:8000")
	http.ListenAndServe(":8000", nil)
}
