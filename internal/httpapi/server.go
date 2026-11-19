package httpapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jamal-gml/hyperion/internal/config"
	"github.com/jamal-gml/hyperion/internal/identity"
	"github.com/jamal-gml/hyperion/internal/mcpserver"
	"github.com/jamal-gml/hyperion/internal/store"
)

type API struct {
	db  *store.Store
	cfg config.Config
	mcp *mcpserver.Server
}

func New(db *store.Store, cfg config.Config) *API {
	return &API{db: db, cfg: cfg, mcp: mcpserver.New(db, identity.Identity{ID: "http-client", Project: "default", Role: "mcp-client", Provider: "mcp"})}
}
func (a *API) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"status": "ok", "service": "hyperion", "time": time.Now().UTC()})
	})
	mux.Handle("/mcp", a.mcp.HTTPHandler())
	mux.HandleFunc("GET /api/agents", a.listAgents)
	mux.HandleFunc("GET /api/tasks", a.listTasks)
	mux.HandleFunc("GET /", a.dashboard)
	srv := &http.Server{Addr: a.cfg.ListenAddr, Handler: securityHeaders(mux), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("Hyperion listening on http://%s", a.cfg.ListenAddr)
	return srv.ListenAndServe()
}
func (a *API) listAgents(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.DB.QueryContext(r.Context(), `SELECT id,project,provider,role,status,last_seen FROM agents ORDER BY last_seen DESC LIMIT 200`)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()
	items := []map[string]string{}
	for rows.Next() {
		var id, p, pr, role, status, seen string
		if rows.Scan(&id, &p, &pr, &role, &status, &seen) != nil {
			continue
		}
		items = append(items, map[string]string{"id": id, "project": p, "provider": pr, "role": role, "status": status, "last_seen": seen})
	}
	writeJSON(w, 200, items)
}
func (a *API) listTasks(w http.ResponseWriter, r *http.Request) {
	project := r.URL.Query().Get("project")
	if project == "" {
		project = "default"
	}
	rows, err := a.db.DB.QueryContext(r.Context(), `SELECT id,title,assigned_to,status,priority,updated_at FROM tasks WHERE project=? ORDER BY updated_at DESC LIMIT 100`, project)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()
	items := []map[string]string{}
	for rows.Next() {
		var id, title, assigned, status, priority, updated string
		if rows.Scan(&id, &title, &assigned, &status, &priority, &updated) != nil {
			continue
		}
		items = append(items, map[string]string{"id": id, "title": title, "assigned_to": assigned, "status": status, "priority": priority, "updated_at": updated})
	}
	writeJSON(w, 200, items)
}
func (a *API) dashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!doctype html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Hyperion</title><style>body{margin:0;background:#0b1020;color:#e8edf8;font:15px Segoe UI,Arial}main{max-width:1100px;margin:48px auto;padding:0 24px}.muted{color:#98a5bd}.grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(260px,1fr));gap:16px}.card{background:#131c31;border:1px solid #263451;border-radius:14px;padding:20px}.pill{display:inline-block;padding:4px 9px;border-radius:99px;background:#123d35;color:#8de0c2;font-size:12px}h1{font-size:32px;margin-bottom:4px}h2{font-size:17px}pre{white-space:pre-wrap;color:#a9b8d2}button{background:#233c67;color:white;border:0;padding:9px 13px;border-radius:8px;cursor:pointer}</style></head><body><main><h1>Hyperion</h1><p class="muted">Local AI collaboration control plane</p><p><span class="pill" id="health">Checking service…</span></p><div class="grid"><section class="card"><h2>Connected / registered agents</h2><div id="agents" class="muted">Loading…</div></section><section class="card"><h2>Recent tasks</h2><div id="tasks" class="muted">Loading…</div></section><section class="card"><h2>MCP endpoint</h2><pre>http://127.0.0.1:8080/mcp</pre><p class="muted">Use the STDIO launcher for separate Claude identities.</p></section><section class="card"><h2>Quick start</h2><pre>scripts\start-agent.ps1 -Project your-project -Role dev -WorkingDirectory C:\work\your-project</pre><p class="muted">Open another terminal with role code-review.</p></section></div><p class="muted">This dashboard shows registered records, not proof that an interactive AI session is currently responsive.</p><button onclick="refresh()">Refresh</button></main><script>async function refresh(){try{let h=await fetch('/healthz');document.getElementById('health').textContent=h.ok?'Service online':'Service unavailable';let a=await (await fetch('/api/agents')).json();document.getElementById('agents').innerHTML=a.length?a.map(x=>'<p><b>'+x.id+'</b><br><span class="muted">'+x.project+' · '+x.provider+' · '+x.role+' · '+x.status+'</span></p>').join(''):'No agents registered yet';let t=await (await fetch('/api/tasks?project=your-project')).json();document.getElementById('tasks').innerHTML=t.length?t.map(x=>'<p><b>'+x.title+'</b><br><span class="muted">'+x.status+' · '+(x.assigned_to||'unassigned')+'</span></p>').join(''):'No your-project tasks yet'}catch(e){document.getElementById('health').textContent='Service unavailable'}}refresh();setInterval(refresh,5000)</script></body></html>`)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}
