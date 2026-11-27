package routes

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"karsa/internal/config"
	"karsa/internal/database"
	"karsa/internal/models"
	"karsa/internal/scheduler"
	"karsa/internal/worker"
)

type api struct {
	test   *testing.T
	router *gin.Engine
	db     *gorm.DB
	suffix string
}

type object = map[string]any

func newAPI(t *testing.T) api {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN belum diisi")
	}
	if !strings.Contains(dsn, "_test?") {
		t.Fatal("nama database di TEST_DB_DSN harus berakhiran _test")
	}
	gin.SetMode(gin.TestMode)
	if err := database.Migrate(dsn, "../../migrations"); err != nil {
		t.Fatal(err)
	}
	db, err := database.Connect(dsn)
	if err != nil {
		t.Fatal(err)
	}
	router, err := New(db, config.Config{JWTSecret: "rahasia", FrontendURL: "http://localhost:5173", APIURL: "http://api", UploadDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	return api{test: t, router: router, db: db, suffix: strings.ToLower(rand.Text()[:10])}
}

func (api api) send(request *http.Request, token string, want int) object {
	api.test.Helper()
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	recorder := httptest.NewRecorder()
	api.router.ServeHTTP(recorder, request)
	if recorder.Code != want {
		api.test.Fatalf("%s %s = %d, want %d: %s", request.Method, request.URL, recorder.Code, want, recorder.Body)
	}
	var result object
	_ = json.Unmarshal(recorder.Body.Bytes(), &result)
	return result
}

func (api api) call(method, path, token string, body any, want int) object {
	api.test.Helper()
	content, _ := json.Marshal(body)
	return api.send(httptest.NewRequest(method, "/api"+path, bytes.NewReader(content)), token, want)
}

func (api api) upload(path, token, field string, content []byte, want int) object {
	api.test.Helper()
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	part, _ := form.CreateFormFile(field, "file.pdf")
	_, _ = part.Write(content)
	_ = form.Close()
	request := httptest.NewRequest("POST", "/api"+path, &body)
	request.Header.Set("Content-Type", form.FormDataContentType())
	return api.send(request, token, want)
}

func (api api) register(name, role string) (string, uint) {
	api.test.Helper()
	result := api.call("POST", "/auth/register", "", object{
		"name": name, "email": name + api.suffix + "@karsa.test", "password": "password123", "role": role, "company_name": "PT " + name + api.suffix,
	}, 201)
	return result["token"].(string), idOf(result["user"])
}

func (api api) admin() string {
	api.test.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	email := "admin" + api.suffix + "@karsa.test"
	if err := api.db.Create(&models.User{Name: "Admin", Email: email, PasswordHash: string(hash), Role: "admin", IsActive: true}).Error; err != nil {
		api.test.Fatal(err)
	}
	return api.call("POST", "/auth/login", "", object{"email": email, "password": "password123"}, 200)["token"].(string)
}

func (api api) verifiedCompany(admin, name string) string {
	api.test.Helper()
	token, _ := api.register(name, "company")
	profile := api.call("GET", "/company/profile", token, nil, 200)
	api.call("POST", fmt.Sprintf("/admin/companies/%d/verify", idOf(profile)), admin, nil, 200)
	return token
}

func (api api) seeker(name string) (string, uint) {
	api.test.Helper()
	token, id := api.register(name, "seeker")
	api.call("PUT", "/me/profile", token, profileBody(false), 200)
	api.upload("/me/cv", token, "cv", []byte("%PDF-1.4 "+name), 200)
	return token, id
}

func (api api) job(admin, company string, changes object) object {
	api.test.Helper()
	category := api.call("POST", "/admin/categories", admin, object{"name": "Kategori " + rand.Text()[:12]}, 201)
	body := object{
		"category_id": category["id"], "title": "Lowongan " + api.suffix, "type": "full_time", "level": "entry", "work_mode": "remote",
		"salary_min": 5000000, "salary_max": 8000000, "show_salary": true, "skills": []string{"Go", "SQL"},
		"description": "Deskripsi", "requirements": "Syarat", "deadline": dateIn(30), "status": "open",
	}
	for key, value := range changes {
		body[key] = value
	}
	return api.call("POST", "/company/jobs", company, body, 201)
}

func profileBody(isDiscoverable bool) object {
	return object{
		"headline": "Backend", "phone": "0812", "location": "Bandung", "education": "S1", "graduation_year": time.Now().Year(),
		"experience_years": 1, "skills": []string{" go ", "GO", "Vue"}, "is_discoverable": isDiscoverable,
	}
}

func dateIn(days int) string {
	return time.Now().AddDate(0, 0, days).Format(time.DateOnly)
}

func idOf(value any) uint {
	return uint(value.(object)["id"].(float64))
}

func (api api) equal(got, want any) {
	api.test.Helper()
	if fmt.Sprint(got) != fmt.Sprint(want) {
		api.test.Fatalf("got %v, want %v", got, want)
	}
}

func TestAuthAndRoles(t *testing.T) {
	api := newAPI(t)
	seeker, seekerID := api.register("sari", "seeker")
	company, _ := api.register("bima", "company")
	admin := api.admin()

	api.call("POST", "/auth/register", "", object{"name": "x", "email": "SARI" + api.suffix + "@karsa.test", "password": "password123", "role": "seeker"}, 422)
	api.call("POST", "/auth/register", "", object{"name": "x", "email": "adm" + api.suffix + "@karsa.test", "password": "password123", "role": "admin"}, 422)
	api.call("POST", "/auth/register", "", object{"name": "x", "email": "c" + api.suffix + "@karsa.test", "password": "password123", "role": "company"}, 422)
	blank := api.call("POST", "/auth/register", "", object{"name": "   ", "email": "salah", "password": "1", "role": "seeker"}, 422)
	api.equal(blank["errors"], object{"name": "Wajib diisi.", "email": "Format email tidak valid.", "password": "Minimal 8."})

	api.call("GET", "/auth/me", "", nil, 401)
	api.call("GET", "/auth/me", "palsu", nil, 401)
	api.call("GET", "/company/jobs", seeker, nil, 403)
	api.call("GET", "/me/profile", company, nil, 403)
	api.call("GET", "/admin/dashboard", company, nil, 403)
	api.call("GET", "/company/talents", company, nil, 403)

	api.call("PATCH", fmt.Sprintf("/admin/users/%d", seekerID), admin, object{"is_active": false}, 200)
	api.equal(api.call("GET", "/auth/me", seeker, nil, 403)["message"], "Akun dinonaktifkan.")
	me := api.call("GET", "/auth/me", admin, nil, 200)
	api.call("PATCH", fmt.Sprintf("/admin/users/%d", idOf(me["user"])), admin, object{"is_active": false}, 422)

	for range 4 {
		api.call("POST", "/auth/login", "", object{"email": "tidakada@karsa.test", "password": "salah"}, 401)
	}
	api.call("POST", "/auth/login", "", object{"email": "tidakada@karsa.test", "password": "salah"}, 429)
}

func TestProfileAndCV(t *testing.T) {
	api := newAPI(t)
	token, _ := api.register("nadia", "seeker")

	api.call("PUT", "/me/profile", token, profileBody(true), 422)
	profile := api.call("PUT", "/me/profile", token, profileBody(false), 200)
	api.equal(profile["skills"], "[go Vue]")
	api.equal(profile["completeness"], object{"percent": 83, "missing": []string{"cv"}})

	api.upload("/me/cv", token, "cv", []byte("bukan pdf"), 422)
	api.upload("/me/cv", token, "cv", append([]byte("%PDF-"), make([]byte, 2<<20)...), 422)
	api.upload("/me/cv", token, "cv", append([]byte("%PDF-"), make([]byte, 5<<20)...), 422)
	api.call("GET", "/me/cv", token, nil, 404)
	profile = api.upload("/me/cv", token, "cv", []byte("%PDF-1.4"), 200)
	api.equal(profile["completeness"].(object)["percent"], 100)
	api.call("GET", "/me/cv", token, nil, 200)

	api.call("PUT", "/me/profile", token, object{"portfolio_url": "javascript:alert(1)"}, 422)
	api.equal(api.call("PUT", "/me/profile", token, profileBody(true), 200)["is_discoverable"], true)
	empty := api.call("PUT", "/me/profile", token, object{"is_discoverable": true}, 200)
	api.equal(empty["is_discoverable"], false)
}

func TestJobRules(t *testing.T) {
	api := newAPI(t)
	admin := api.admin()
	company := api.verifiedCompany(admin, "nusa")
	other := api.verifiedCompany(admin, "lain")
	pending, _ := api.register("baru", "company")

	job := api.job(admin, company, object{"screening_questions": []object{{"question": "Siap?", "type": "yes_no", "required": true, "expected_answer": "yes"}}})
	path := fmt.Sprintf("/company/jobs/%d", idOf(job))
	body := object{}
	for _, key := range []string{"category_id", "title", "type", "level", "work_mode", "salary_min", "salary_max", "skills", "description", "requirements", "deadline"} {
		body[key] = job[key]
	}
	with := func(changes object) object {
		merged := object{"status": "open"}
		for key, value := range body {
			merged[key] = value
		}
		for key, value := range changes {
			merged[key] = value
		}
		return merged
	}

	api.call("POST", "/company/jobs", pending, with(nil), 403)
	api.call("POST", "/company/jobs", pending, with(object{"status": "draft"}), 201)
	api.call("GET", path, other, nil, 404)
	api.call("DELETE", path, other, nil, 404)
	api.call("PUT", path, company, with(object{"status": "draft"}), 422)
	api.call("PUT", path, company, with(object{"deadline": dateIn(-1)}), 422)
	api.call("PUT", path, company, with(object{"deadline": dateIn(91)}), 422)
	api.call("PUT", path, company, with(object{"work_mode": "onsite", "location": " "}), 422)
	api.call("PUT", path, company, with(object{"salary_max": 1}), 422)
	invalid := api.call("PUT", path, company, with(object{"screening_questions": []object{{"question": "Umur?", "type": "number", "expected_answer": ">5"}}}), 422)
	api.equal(len(invalid["errors"].(object)), 1)

	detail := api.call("GET", "/jobs/"+job["slug"].(string), "", nil, 200)
	question := detail["screening_questions"].([]any)[0].(object)
	if _, leaked := question["expected_answer"]; leaked {
		t.Fatal("expected_answer bocor ke publik")
	}
	hidden := api.call("PUT", path, company, with(object{"show_salary": false}), 200)
	api.equal(hidden["salary_min"], 5e+06)
	api.equal(api.call("GET", "/jobs/"+job["slug"].(string), "", nil, 200)["salary_min"], nil)

	api.call("POST", path+"/close", company, nil, 200)
	api.call("GET", "/jobs/"+job["slug"].(string), "", nil, 404)
	api.call("POST", path+"/open", company, nil, 200)
	api.call("POST", fmt.Sprintf("/admin/jobs/%d/close", idOf(job)), admin, object{"reason": "Melanggar"}, 200)
	api.call("POST", path+"/open", company, nil, 403)
	api.call("PUT", path, company, with(nil), 403)
	api.call("DELETE", fmt.Sprintf("/admin/categories/%v", job["category_id"]), admin, nil, 409)
	api.call("DELETE", path, company, nil, 204)
	unused := api.call("POST", "/admin/categories", admin, object{"name": "Kategori " + api.suffix}, 201)
	api.call("POST", "/admin/categories", admin, object{"name": "Kategori " + api.suffix}, 422)
	api.call("DELETE", fmt.Sprintf("/admin/categories/%d", idOf(unused)), admin, nil, 204)
}

func TestSearch(t *testing.T) {
	api := newAPI(t)
	admin := api.admin()
	company := api.verifiedCompany(admin, "cari")
	word := "zx" + api.suffix
	api.job(admin, company, object{"title": word + " pertama", "level": "internship", "salary_max": 9000000})
	api.job(admin, company, object{"title": "Kedua", "skills": []string{word}, "salary_max": 6000000})
	api.job(admin, company, object{"title": "Ketiga", "description": "tentang " + word, "work_mode": "onsite", "location": "Kota " + word, "show_salary": false})
	api.job(admin, company, object{"title": word + " draft", "status": "draft"})

	total := func(query string) any {
		t.Helper()
		return api.call("GET", "/jobs?"+query, "", nil, 200)["total"]
	}
	api.equal(total("q="+word), 3)
	api.equal(total("q="+word+"&level=internship,mid"), 1)
	api.equal(total("q="+word+"&mode=onsite&location="+word), 1)
	api.equal(total("q="+word+"&salary_min=7000000"), 1)
	api.equal(total("q=PT+cari"+api.suffix), 3)
	for _, special := range []string{"%25", "'", "%22%2B-*()", "%5C", "c%2B%2B"} {
		total("q=" + special)
	}
	api.equal(total("q=%25"+api.suffix), 0)
	sorted := api.call("GET", "/jobs?q="+word+"&sort=salary", "", nil, 200)["data"].([]any)
	api.equal(sorted[0].(object)["title"], word+" pertama")
	api.equal(sorted[2].(object)["title"], "Ketiga")
	api.equal(api.call("GET", "/jobs?q="+word+"&page=999999999999", "", nil, 200)["page"], 1)
	api.call("GET", "/home", "", nil, 200)
}

func TestSavedJobs(t *testing.T) {
	api := newAPI(t)
	admin := api.admin()
	company := api.verifiedCompany(admin, "simpan")
	seeker, _ := api.seeker("tika")
	job := api.job(admin, company, nil)
	path := fmt.Sprintf("/me/saved-jobs/%d", idOf(job))

	api.call("PUT", path, seeker, nil, 204)
	api.call("PUT", path, seeker, nil, 204)
	api.equal(api.call("GET", "/jobs/"+job["slug"].(string), seeker, nil, 200)["is_saved"], true)
	api.call("POST", fmt.Sprintf("/company/jobs/%d/close", idOf(job)), company, nil, 200)
	saved := api.call("GET", "/me/saved-jobs", seeker, nil, 200)
	api.equal(saved["total"], 1)
	api.equal(saved["data"].([]any)[0].(object)["is_open"], false)
	api.call("DELETE", path, seeker, nil, 204)
	api.call("DELETE", path, seeker, nil, 204)
}

func TestApplyAndPipeline(t *testing.T) {
	api := newAPI(t)
	admin := api.admin()
	company := api.verifiedCompany(admin, "rekrut")
	other := api.verifiedCompany(admin, "asing")
	seeker, _ := api.seeker("budi")
	stranger, _ := api.seeker("orang")
	incomplete, _ := api.register("kosong", "seeker")
	job := api.job(admin, company, object{"screening_questions": []object{
		{"question": "Siap WFO?", "type": "yes_no", "required": true, "expected_answer": "yes"},
		{"question": "Pengalaman?", "type": "number", "required": false, "expected_answer": ">=1"},
	}})
	questions := job["screening_questions"].([]any)
	applyPath := fmt.Sprintf("/jobs/%d/apply", idOf(job))
	answers := func(first, second string) object {
		return object{"answers": []object{{"question_id": idOf(questions[0]), "answer": first}, {"question_id": idOf(questions[1]), "answer": second}}}
	}

	api.call("POST", applyPath, incomplete, object{}, 422)
	api.equal(len(api.call("POST", applyPath, seeker, answers("", "banyak"), 422)["errors"].(object)), 2)

	var waiting sync.WaitGroup
	codes := make(chan int, 5)
	for range 5 {
		waiting.Go(func() {
			content, _ := json.Marshal(answers("no", "2"))
			request := httptest.NewRequest("POST", "/api"+applyPath, bytes.NewReader(content))
			request.Header.Set("Authorization", "Bearer "+seeker)
			recorder := httptest.NewRecorder()
			api.router.ServeHTTP(recorder, request)
			codes <- recorder.Code
		})
	}
	waiting.Wait()
	close(codes)
	counts := map[int]int{}
	for code := range codes {
		counts[code]++
	}
	api.equal(counts, map[int]int{201: 1, 409: 4})

	api.call("DELETE", fmt.Sprintf("/company/jobs/%d", idOf(job)), company, nil, 409)
	list := api.call("GET", fmt.Sprintf("/company/jobs/%d/applications?min_score=60&skill=GO&meets=false&source=applied", idOf(job)), company, nil, 200)
	api.equal(list["total"], 1)
	api.equal(list["counts"].(object)["applied"], 1)
	row := list["data"].([]any)[0].(object)
	api.equal(row["match_score"], 60)
	path := fmt.Sprintf("/company/applications/%d", idOf(row))
	mine := fmt.Sprintf("/me/applications/%d", idOf(row))

	api.call("GET", path, other, nil, 404)
	api.call("GET", path+"/cv", other, nil, 404)
	api.call("PATCH", path, other, object{"status": "reviewed"}, 404)
	api.call("GET", mine, stranger, nil, 404)
	api.call("GET", mine+"/cv", stranger, nil, 404)

	api.upload("/me/cv", seeker, "cv", []byte("%PDF-1.4 baru"), 200)
	request := httptest.NewRequest("GET", "/api"+path+"/cv", nil)
	request.Header.Set("Authorization", "Bearer "+company)
	recorder := httptest.NewRecorder()
	api.router.ServeHTTP(recorder, request)
	api.equal(recorder.Body.String(), "%PDF-1.4 budi")

	api.call("PATCH", path, company, object{"status": "accepted"}, 422)
	api.equal(api.call("PATCH", path, company, object{"company_note": "Rahasia"}, 200)["company_note"], "Rahasia")
	for _, status := range []string{"reviewed", "interview", "accepted"} {
		api.call("PATCH", path, company, object{"status": status}, 200)
	}
	api.call("PATCH", path, company, object{"status": "rejected"}, 422)
	api.call("POST", mine+"/withdraw", seeker, nil, 422)
	api.equal(api.call("GET", "/me/notifications", seeker, nil, 200)["applications"], 1)
	api.equal(api.call("GET", "/me/applications", seeker, nil, 200)["data"].([]any)[0].(object)["has_update"], true)
	detail := api.call("GET", mine, seeker, nil, 200)
	api.equal(api.call("GET", "/me/notifications", seeker, nil, 200)["applications"], 0)
	api.equal(len(detail["events"].([]any)), 4)
	if _, leaked := detail["company_note"]; leaked {
		t.Fatal("company_note bocor ke pelamar")
	}
	var emails int64
	api.db.Model(&models.EmailJob{}).Where("to_email = ?", "budi"+api.suffix+"@karsa.test").Count(&emails)
	api.equal(emails, 4)

	api.call("POST", applyPath, stranger, answers("yes", "2"), 201)
	applied := api.call("GET", "/me/applications", stranger, nil, 200)["data"].([]any)[0]
	strangerPath := fmt.Sprintf("/me/applications/%d", idOf(applied))
	api.equal(api.call("POST", strangerPath+"/withdraw", stranger, nil, 200)["status"], "withdrawn")
	api.call("POST", applyPath, stranger, answers("yes", "2"), 409)

	dashboard := api.call("GET", "/company/dashboard", company, nil, 200)
	api.equal(dashboard["total_applicants"], 2)
	api.equal(len(dashboard["weekly"].([]any)), 8)
	api.equal(dashboard["weekly"].([]any)[7].(object)["count"], 2)
	api.call("GET", "/admin/dashboard", admin, nil, 200)
	api.call("GET", "/admin/jobs?status=open&company=rekrut", admin, nil, 200)
	api.call("GET", "/admin/users?q="+api.suffix+"&role=seeker", admin, nil, 200)
	api.call("GET", "/admin/companies?status=verified", admin, nil, 200)
}

func TestConcurrentStatusChange(t *testing.T) {
	api := newAPI(t)
	admin := api.admin()
	company := api.verifiedCompany(admin, "balap")
	seeker, _ := api.seeker("cepat")
	job := api.job(admin, company, nil)
	application := api.call("POST", fmt.Sprintf("/jobs/%d/apply", idOf(job)), seeker, object{}, 201)["application"]

	var waiting sync.WaitGroup
	codes := make(chan int, 6)
	for range 6 {
		waiting.Go(func() {
			request := httptest.NewRequest("PATCH", fmt.Sprintf("/api/company/applications/%d", idOf(application)), strings.NewReader(`{"status":"reviewed"}`))
			request.Header.Set("Authorization", "Bearer "+company)
			recorder := httptest.NewRecorder()
			api.router.ServeHTTP(recorder, request)
			codes <- recorder.Code
		})
	}
	waiting.Wait()
	close(codes)
	succeeded := 0
	for code := range codes {
		if code == 200 {
			succeeded++
		} else if code != 409 && code != 422 {
			t.Fatalf("kode tidak terduga: %d", code)
		}
	}
	api.equal(succeeded, 1)
}

func TestTalentsAndInvitations(t *testing.T) {
	api := newAPI(t)
	admin := api.admin()
	company := api.verifiedCompany(admin, "undang")
	seeker, seekerID := api.seeker("talenta")
	_, hiddenID := api.seeker("tertutup")
	skill := "sk" + api.suffix
	profile := profileBody(true)
	profile["skills"] = []string{skill, "Vue"}
	api.call("PUT", "/me/profile", seeker, profile, 200)
	job := api.job(admin, company, object{"skills": []string{skill}})
	second := api.job(admin, company, nil)

	talents := api.call("GET", "/company/talents?max_experience=1&education=S1&skills=vue,"+strings.ToUpper(skill), company, nil, 200)
	api.equal(talents["total"], 1)
	talent := talents["data"].([]any)[0].(object)
	for _, secret := range []string{"email", "phone", "cv", "cv_url", "cv_path", "bio", "portfolio_url"} {
		if _, leaked := talent[secret]; leaked {
			t.Fatalf("%s bocor di Cari Talenta", secret)
		}
	}

	invite := object{"user_id": seekerID, "job_id": idOf(job), "message": "Halo"}
	api.call("POST", "/company/invitations", company, object{"user_id": hiddenID, "job_id": idOf(job)}, 404)
	invitation := api.call("POST", "/company/invitations", company, invite, 201)["invitation"]
	api.call("POST", "/company/invitations", company, invite, 409)
	declined := api.call("POST", "/company/invitations", company, object{"user_id": seekerID, "job_id": idOf(second)}, 201)["invitation"]
	api.equal(api.call("GET", "/company/invitations", company, nil, 200)["remaining_today"], 18)

	stranger, _ := api.seeker("penyusup")
	applyPath := fmt.Sprintf("/jobs/%d/apply", idOf(job))
	api.call("POST", applyPath, stranger, object{"invitation_id": idOf(invitation)}, 422)
	api.call("POST", fmt.Sprintf("/me/invitations/%d/decline", idOf(invitation)), stranger, nil, 404)
	api.equal(api.call("GET", "/me/notifications", seeker, nil, 200)["invitations"], 2)
	application := api.call("POST", applyPath, seeker, object{"invitation_id": idOf(invitation)}, 201)["application"].(object)
	api.equal(application["source"], "invited")
	api.equal(application["status"], "applied")
	api.equal(api.call("GET", "/me/notifications", seeker, nil, 200)["invitations"], 1)
	api.call("POST", fmt.Sprintf("/me/invitations/%d/decline", idOf(declined)), seeker, nil, 200)
	api.call("POST", fmt.Sprintf("/me/invitations/%d/decline", idOf(declined)), seeker, nil, 422)
	api.equal(api.call("GET", "/company/invitations?status=accepted", company, nil, 200)["total"], 1)
	applicant := api.call("GET", fmt.Sprintf("/company/applications/%d", idOf(application)), company, nil, 200)
	api.equal(applicant["candidate"].(object)["email"], "talenta"+api.suffix+"@karsa.test")

	var owner models.Invitation
	api.db.First(&owner, idOf(invitation))
	for range 18 {
		filler := api.job(admin, company, nil)
		api.db.Create(&models.Invitation{CompanyID: owner.CompanyID, JobID: idOf(filler), UserID: seekerID, Status: "pending"})
	}
	last := api.job(admin, company, nil)
	api.call("POST", "/company/invitations", company, object{"user_id": seekerID, "job_id": idOf(last)}, 429)
}

func TestRejectedCompanyResubmits(t *testing.T) {
	api := newAPI(t)
	admin := api.admin()
	company, _ := api.register("ulang", "company")
	id := idOf(api.call("GET", "/company/profile", company, nil, 200))
	path := fmt.Sprintf("/admin/companies/%d", id)

	api.call("POST", path+"/reject", admin, object{"reason": " "}, 422)
	api.call("POST", path+"/reject", admin, object{"reason": "Profil kosong"}, 200)
	api.call("POST", path+"/reject", admin, object{"reason": "Profil kosong"}, 409)
	api.equal(api.call("GET", "/company/profile", company, nil, 200)["rejection_reason"], "Profil kosong")
	profile := api.call("PUT", "/company/profile", company, object{"name": "PT Baru", "size": ">500", "website": "https://baru.example"}, 200)
	api.equal(profile["status"], "pending")
	api.equal(profile["rejection_reason"], nil)
	api.upload("/company/logo", company, "logo", []byte("<svg onload=alert(1)>"), 422)
	logo := api.upload("/company/logo", company, "logo", []byte("\x89PNG\r\n\x1a\n"), 200)
	if !strings.HasPrefix(fmt.Sprint(logo["logo_url"]), "http://api/uploads/logos/") {
		t.Fatalf("logo_url = %v", logo["logo_url"])
	}
	logoPath := strings.TrimPrefix(logo["logo_url"].(string), "http://api")
	api.send(httptest.NewRequest("GET", logoPath, nil), "", 200)
	api.send(httptest.NewRequest("GET", "/uploads/cv/", nil), "", 404)
	api.call("POST", path+"/verify", admin, nil, 200)
}

func TestSchedulerClosesAndExpires(t *testing.T) {
	api := newAPI(t)
	admin := api.admin()
	company := api.verifiedCompany(admin, "jadwal")
	_, seekerID := api.seeker("lama")
	expired := api.job(admin, company, nil)
	active := api.job(admin, company, nil)
	var owner models.Job
	api.db.First(&owner, idOf(expired))
	api.db.Model(&owner).Update("deadline", dateIn(-1))
	invitations := []models.Invitation{
		{CompanyID: owner.CompanyID, JobID: idOf(expired), UserID: seekerID, Status: "pending"},
		{CompanyID: owner.CompanyID, JobID: idOf(active), UserID: seekerID, Status: "pending", CreatedAt: time.Now().AddDate(0, 0, -15)},
	}
	api.db.Create(&invitations)

	for range 2 {
		if err := scheduler.Sweep(api.db); err != nil {
			t.Fatal(err)
		}
	}
	api.equal(api.call("GET", fmt.Sprintf("/company/jobs/%d", idOf(expired)), company, nil, 200)["status"], "closed")
	api.equal(api.call("GET", fmt.Sprintf("/company/jobs/%d", idOf(active)), company, nil, 200)["status"], "open")
	api.equal(api.call("GET", "/company/invitations?status=expired", company, nil, 200)["total"], 2)
}

type fakeSender struct {
	mutex sync.Mutex
	fail  bool
	sent  map[uint]int
}

func (sender *fakeSender) Send(job models.EmailJob) error {
	sender.mutex.Lock()
	defer sender.mutex.Unlock()
	if sender.fail {
		return errors.New("smtp mati")
	}
	sender.sent[job.ID]++
	return nil
}

func TestEmailWorker(t *testing.T) {
	api := newAPI(t)
	api.db.Where("1 = 1").Delete(&models.EmailJob{})
	jobs := make([]models.EmailJob, 12)
	for index := range jobs {
		jobs[index] = models.EmailJob{ToEmail: "a@karsa.test", Status: "pending", AvailableAt: time.Now().Add(-time.Second)}
	}
	api.db.Create(&jobs)

	sender := &fakeSender{sent: map[uint]int{}}
	var waiting sync.WaitGroup
	for range 4 {
		waiting.Go(func() {
			for {
				if processed, err := worker.ProcessNext(api.db, sender); err != nil || !processed {
					return
				}
			}
		})
	}
	waiting.Wait()
	api.equal(len(sender.sent), 12)
	for id, count := range sender.sent {
		if count != 1 {
			t.Fatalf("email %d terkirim %d kali", id, count)
		}
	}

	failing := models.EmailJob{ToEmail: "b@karsa.test", Status: "pending", AvailableAt: time.Now().Add(-time.Second)}
	api.db.Create(&failing)
	sender.fail = true
	for attempt := 1; attempt <= 3; attempt++ {
		if processed, err := worker.ProcessNext(api.db, sender); err != nil || !processed {
			t.Fatalf("percobaan %d: %v %v", attempt, processed, err)
		}
		if processed, _ := worker.ProcessNext(api.db, sender); processed {
			t.Fatal("email diulang sebelum waktunya")
		}
		api.db.First(&failing, failing.ID)
		api.equal(failing.Attempts, attempt)
		api.db.Model(&failing).Update("available_at", time.Now().Add(-time.Second))
	}
	api.equal(failing.Status, "failed")
}

func TestAvatarUpload(t *testing.T) {
	api := newAPI(t)
	token, _ := api.register("foto", "seeker")
	api.upload("/me/avatar", token, "avatar", []byte("<svg onload=alert(1)>"), 422)
	api.upload("/me/avatar", token, "avatar", append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 3<<20)...), 422)
	first := api.upload("/me/avatar", token, "avatar", []byte("\x89PNG\r\n\x1a\nsatu"), 200)["user"].(object)["avatar_url"].(string)
	if !strings.HasPrefix(first, "http://api/uploads/avatars/") {
		t.Fatalf("avatar_url = %s", first)
	}
	api.send(httptest.NewRequest("GET", strings.TrimPrefix(first, "http://api"), nil), "", 200)
	second := api.upload("/me/avatar", token, "avatar", []byte("\x89PNG\r\n\x1a\ndua"), 200)["user"].(object)["avatar_url"].(string)
	api.send(httptest.NewRequest("GET", strings.TrimPrefix(first, "http://api"), nil), "", 404)
	api.equal(api.call("GET", "/auth/me", token, nil, 200)["user"].(object)["avatar_url"], second)
}
