package seed

type companyData struct {
	name, industry, size, location, website, description, status string
}

type jobData struct {
	company                                         int
	category, title, jobType, level, mode, location string
	salaryMin, salaryMax                            uint
	skills                                          []string
	deadlineDays                                    int
	status                                          string
}

type seekerData struct {
	name, email, headline, location, education, major string
	graduationOffset, experience                      int
	skills                                            []string
	isDiscoverable                                    bool
}

type applicationData struct {
	seeker, job int
	statuses    []string
	daysAgo     int
}

var categoryNames = []string{"Software Engineering", "Data", "UI/UX Design", "Product", "QA", "DevOps & Cloud", "Digital Marketing"}

var companies = []companyData{
	{"Nusa Digital", "Software House", "51-200", "Jakarta", "https://nusadigital.example", "Membangun aplikasi web dan mobile untuk perusahaan di Indonesia.", "verified"},
	{"Laras Data", "Data & Analitik", "11-50", "Bandung", "https://larasdata.example", "Platform analitik untuk bisnis ritel skala menengah.", "verified"},
	{"Kembang Studio", "Desain Produk", "11-50", "Yogyakarta", "https://kembang.example", "Studio desain produk digital dengan tim kecil dan riset yang kuat.", "verified"},
	{"Awan Teknologi", "Cloud", "201-500", "Surabaya", "https://awantek.example", "Penyedia layanan cloud dan infrastruktur terkelola.", "verified"},
	{"Pasar Pintar", "E-commerce", "51-200", "Jakarta", "https://pasarpintar.example", "Marketplace kebutuhan harian untuk kota tingkat dua.", "verified"},
	{"Belajar Yuk", "Edukasi", "11-50", "Semarang", "https://belajaryuk.example", "Aplikasi belajar interaktif untuk siswa SMA.", "verified"},
	{"Sinar Muda Labs", "Startup", "1-10", "Medan", "https://sinarmuda.example", "Startup baru yang sedang membangun produk pertamanya.", "pending"},
}

var jobs = []jobData{
	{0, "Software Engineering", "Frontend Developer Intern", "internship", "internship", "hybrid", "Jakarta", 2, 3, []string{"Vue", "JavaScript", "CSS"}, 40, "open"},
	{0, "Software Engineering", "Junior Backend Developer (Go)", "full_time", "junior", "hybrid", "Jakarta", 7, 11, []string{"Go", "MySQL", "Docker"}, 50, "open"},
	{0, "Software Engineering", "Fullstack Developer", "full_time", "mid", "remote", "", 12, 18, []string{"Go", "Vue", "MySQL"}, 60, "open"},
	{0, "QA", "QA Engineer Entry Level", "full_time", "entry", "onsite", "Jakarta", 5, 7, []string{"Testing", "Postman"}, 30, "open"},
	{1, "Data", "Data Analyst Intern", "internship", "internship", "onsite", "Bandung", 2, 3, []string{"SQL", "Excel", "Python"}, 25, "open"},
	{1, "Data", "Junior Data Analyst", "full_time", "junior", "hybrid", "Bandung", 6, 9, []string{"SQL", "Python", "Tableau"}, 45, "open"},
	{1, "Data", "Data Engineer", "full_time", "mid", "remote", "", 13, 20, []string{"Python", "SQL", "Airflow"}, 70, "open"},
	{1, "Software Engineering", "Backend Developer Python", "full_time", "entry", "hybrid", "Bandung", 6, 8, []string{"Python", "SQL"}, 35, "open"},
	{2, "UI/UX Design", "UI/UX Designer Intern", "internship", "internship", "remote", "", 2, 3, []string{"Figma", "Wireframing"}, 20, "open"},
	{2, "UI/UX Design", "Junior Product Designer", "full_time", "junior", "onsite", "Yogyakarta", 5, 8, []string{"Figma", "Prototyping", "User Research"}, 40, "open"},
	{2, "Product", "Associate Product Manager", "full_time", "entry", "hybrid", "Yogyakarta", 6, 9, []string{"Roadmap", "User Research"}, 55, "open"},
	{2, "Digital Marketing", "Content Marketing Intern", "internship", "internship", "remote", "", 1, 2, []string{"Copywriting", "SEO"}, 15, "open"},
	{3, "DevOps & Cloud", "Cloud Support Intern", "internship", "internship", "onsite", "Surabaya", 2, 3, []string{"Linux", "AWS"}, 30, "open"},
	{3, "DevOps & Cloud", "Junior DevOps Engineer", "full_time", "junior", "hybrid", "Surabaya", 8, 12, []string{"Linux", "Docker", "CI/CD"}, 60, "open"},
	{3, "DevOps & Cloud", "Site Reliability Engineer", "full_time", "mid", "remote", "", 14, 22, []string{"Kubernetes", "Go", "Linux"}, 80, "open"},
	{3, "QA", "QA Automation Intern", "internship", "internship", "hybrid", "Surabaya", 2, 3, []string{"Testing", "JavaScript"}, 28, "open"},
	{4, "Software Engineering", "Junior Frontend Developer (Vue)", "full_time", "junior", "onsite", "Jakarta", 6, 10, []string{"Vue", "TypeScript", "CSS"}, 42, "open"},
	{4, "Software Engineering", "Mobile Developer Kontrak", "contract", "junior", "remote", "", 8, 12, []string{"Flutter", "REST API"}, 33, "open"},
	{4, "Digital Marketing", "Performance Marketing Entry Level", "full_time", "entry", "hybrid", "Jakarta", 5, 7, []string{"Google Ads", "Excel"}, 38, "open"},
	{4, "Product", "Product Analyst Paruh Waktu", "part_time", "entry", "remote", "", 3, 5, []string{"SQL", "Excel"}, 26, "open"},
	{5, "Software Engineering", "Web Developer Intern", "internship", "internship", "remote", "", 2, 3, []string{"HTML", "CSS", "JavaScript"}, 18, "open"},
	{5, "UI/UX Design", "Visual Designer Paruh Waktu", "part_time", "entry", "remote", "", 2, 4, []string{"Figma", "Illustrator"}, 22, "open"},
	{5, "QA", "QA Tester Intern", "internship", "internship", "onsite", "Semarang", 1, 2, []string{"Testing"}, 2, "open"},
	{5, "Digital Marketing", "Social Media Specialist", "full_time", "entry", "onsite", "Semarang", 4, 6, []string{"Copywriting", "Canva"}, 34, "closed"},
	{6, "Software Engineering", "Founding Engineer", "full_time", "mid", "remote", "", 10, 15, []string{"Go", "Vue"}, 45, "draft"},
	{0, "Software Engineering", "Backend Developer Intern (Node.js)", "internship", "internship", "hybrid", "Jakarta", 2, 3, []string{"Node.js", "Express", "MongoDB"}, 35, "open"},
	{0, "Software Engineering", "Android Developer", "full_time", "junior", "onsite", "Jakarta", 8, 12, []string{"Kotlin", "Android"}, 50, "open"},
	{0, "Software Engineering", "iOS Developer Intern", "internship", "internship", "remote", "", 2, 4, []string{"Swift", "iOS"}, 28, "open"},
	{0, "Data", "Data Engineer Intern", "internship", "internship", "hybrid", "Jakarta", 2, 3, []string{"Python", "SQL"}, 21, "open"},
	{0, "Data", "Machine Learning Engineer", "full_time", "mid", "remote", "", 15, 24, []string{"Python", "TensorFlow", "SQL"}, 75, "open"},
	{0, "Data", "Business Intelligence Analyst", "full_time", "entry", "onsite", "Jakarta", 6, 8, []string{"SQL", "Power BI"}, 3, "open"},
	{0, "UI/UX Design", "UI Designer", "full_time", "junior", "hybrid", "Jakarta", 7, 10, []string{"Figma", "Design System"}, 44, "open"},
	{0, "UI/UX Design", "UX Researcher Intern", "internship", "internship", "remote", "", 2, 3, []string{"User Research", "Figma"}, 30, "open"},
	{0, "Product", "Product Manager", "full_time", "mid", "hybrid", "Jakarta", 16, 25, []string{"Roadmap", "Agile", "SQL"}, 60, "open"},
	{0, "Product", "Scrum Master", "full_time", "junior", "onsite", "Jakarta", 9, 13, []string{"Agile", "Jira"}, 38, "open"},
	{0, "QA", "QA Automation Engineer", "full_time", "junior", "hybrid", "Jakarta", 8, 12, []string{"Selenium", "JavaScript", "Testing"}, 47, "open"},
	{0, "QA", "Manual Tester Paruh Waktu", "part_time", "entry", "remote", "", 3, 4, []string{"Testing", "Jira"}, 25, "open"},
	{0, "DevOps & Cloud", "Cloud Engineer", "full_time", "mid", "remote", "", 15, 22, []string{"AWS", "Terraform", "Linux"}, 70, "open"},
	{0, "DevOps & Cloud", "DevOps Intern", "internship", "internship", "onsite", "Jakarta", 2, 3, []string{"Docker", "Linux", "CI/CD"}, 18, "open"},
	{0, "Digital Marketing", "SEO Specialist", "full_time", "entry", "hybrid", "Jakarta", 5, 7, []string{"SEO", "Google Analytics"}, 33, "open"},
	{0, "Digital Marketing", "Social Media Intern", "internship", "internship", "onsite", "Jakarta", 1, 2, []string{"Canva", "Copywriting"}, 40, "draft"},
	{0, "Software Engineering", "Laravel Developer Kontrak", "contract", "junior", "remote", "", 8, 11, []string{"PHP", "Laravel", "MySQL"}, 52, "draft"},
	{0, "Software Engineering", "React Native Developer", "full_time", "junior", "hybrid", "Jakarta", 9, 13, []string{"React Native", "JavaScript"}, 36, "closed"},
}

var seekers = []seekerData{
	{"Nadia Putri", "nadia@karsa.test", "Mahasiswa Informatika, fokus frontend", "Jakarta", "S1", "Informatika", 1, 0, []string{"Vue", "JavaScript", "CSS", "Figma"}, true},
	{"Raka Pratama", "raka@karsa.test", "Backend developer junior", "Bandung", "S1", "Sistem Informasi", -2, 2, []string{"Go", "MySQL", "Docker", "Python"}, false},
	{"Dinda Larasati", "dinda@karsa.test", "Fresh graduate, tertarik data", "Bandung", "S1", "Statistika", 0, 0, []string{"SQL", "Python", "Excel"}, false},
	{"Fajar Nugroho", "fajar@karsa.test", "Calon QA engineer", "Jakarta", "D3", "Teknik Komputer", 0, 0, []string{"Testing", "Postman", "JavaScript"}, true},
	{"Salsa Ramadhani", "salsa@karsa.test", "UI/UX designer pemula", "Yogyakarta", "S1", "Desain Komunikasi Visual", 0, 1, []string{"Figma", "Wireframing", "User Research"}, true},
	{"Galih Saputra", "galih@karsa.test", "DevOps junior", "Surabaya", "S1", "Teknik Informatika", -3, 3, []string{"Linux", "Docker", "CI/CD", "AWS"}, true},
	{"Maya Anggraini", "maya@karsa.test", "Content writer dan SEO", "Semarang", "S1", "Ilmu Komunikasi", -1, 1, []string{"Copywriting", "SEO", "Canva"}, true},
	{"Yoga Permana", "yoga@karsa.test", "Fullstack developer", "Medan", "S1", "Ilmu Komputer", -4, 4, []string{"Go", "Vue", "MySQL", "TypeScript"}, true},
	{"Intan Maharani", "intan@karsa.test", "Data analyst junior", "Jakarta", "S1", "Matematika", -2, 2, []string{"SQL", "Tableau", "Python"}, true},
	{"Bagus Wijaya", "bagus@karsa.test", "Siswa SMK RPL, siap magang", "Semarang", "SMA/SMK", "Rekayasa Perangkat Lunak", 0, 0, []string{"HTML", "CSS", "JavaScript"}, true},
	{"Citra Lestari", "citra@karsa.test", "Associate product manager", "Yogyakarta", "S2", "Manajemen", -1, 1, []string{"Roadmap", "User Research", "SQL"}, true},
}

var applications = []applicationData{
	{0, 0, []string{"reviewed", "interview"}, 12},
	{0, 16, []string{"reviewed"}, 6},
	{0, 8, nil, 1},
	{0, 20, []string{"rejected"}, 20},
	{1, 1, []string{"reviewed", "interview", "accepted"}, 35},
	{1, 7, []string{"reviewed"}, 9},
	{1, 13, nil, 3},
	{2, 4, []string{"reviewed", "interview"}, 15},
	{2, 5, nil, 4},
	{2, 19, []string{"withdrawn"}, 25},
	{3, 3, []string{"reviewed"}, 18},
	{3, 15, nil, 2},
	{4, 8, []string{"reviewed", "interview", "accepted"}, 40},
	{4, 9, nil, 5},
	{5, 13, []string{"reviewed", "interview"}, 28},
	{6, 11, []string{"reviewed", "rejected"}, 45},
	{7, 1, []string{"reviewed"}, 8},
	{7, 2, nil, 2},
	{8, 5, []string{"reviewed", "interview"}, 22},
	{9, 20, nil, 11},
	{9, 0, nil, 3},
	{10, 10, []string{"reviewed"}, 14},
}
