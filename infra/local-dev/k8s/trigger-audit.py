import urllib.request, json, urllib.error, time

BASE = "http://localhost:8000/api/v1"

def post(path, body, token=None):
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    req = urllib.request.Request(f"{BASE}{path}", data=json.dumps(body).encode(), headers=headers, method="POST")
    try:
        r = urllib.request.urlopen(req, timeout=10)
        return r.status, r.read().decode()
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode()

# 1. Login (should record user.login)
st, body = post("/auth/login", {"username": "admin", "password": "Admin@123"})
print("LOGIN", st)
token = json.loads(body).get("access_token") if st == 200 else None
print("TOKEN", "yes" if token else "no")

if token:
    # 2. List connections to find the seeded connection
    req = urllib.request.Request(f"{BASE}/connections", headers={"Authorization": f"Bearer {token}"})
    try:
        r = urllib.request.urlopen(req, timeout=10)
        conns = json.loads(r.read().decode())
        print("CONNECTIONS", len(conns) if isinstance(conns, list) else conns)
        if isinstance(conns, list) and len(conns) > 0:
            cid = conns[0].get("connection_id")
            print("CONN_ID", cid)
            # 3. Trigger sync (should record connection.sync, connection_run.start, connection_run.complete)
            st2, body2 = post(f"/connections/{cid}/sync", {"sync_type": "INCREMENTAL"}, token)
            print("SYNC", st2, body2[:200])
    except Exception as e:
        print("CONNS_ERR", e)

# 4. Check audit_logs
import psycopg2
conn = psycopg2.connect("postgresql://fusion:fusion_local@postgres:5432/fusion_cdc_metadata")
cur = conn.cursor()
cur.execute("SELECT action, COUNT(*) FROM audit_logs GROUP BY action ORDER BY action")
rows = cur.fetchall()
print("AUDIT_LOGS", rows)
cur.close(); conn.close()
