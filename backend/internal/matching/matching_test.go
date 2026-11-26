package matching

import (
	"slices"
	"testing"

	"karsa/internal/models"
)

func pointer[Value any](value Value) *Value { return &value }

func TestScore(t *testing.T) {
	job := models.Job{Level: "entry", Skills: []string{"Go", "SQL", "Docker"}}
	profile := models.SeekerProfile{Skills: []string{"go", "sql"}, ExperienceYears: 1}

	total, breakdown := Score(job, profile, []bool{true, false}, 2026)
	if total != 70 {
		t.Fatalf("total = %d, want 70", total)
	}
	if !slices.Equal(breakdown.Skill.Matched, []string{"Go", "SQL"}) || !slices.Equal(breakdown.Skill.Missing, []string{"Docker"}) {
		t.Fatalf("skill = %+v", breakdown.Skill)
	}
	if breakdown.Screening.Passed != 1 || breakdown.Screening.Total != 2 {
		t.Fatalf("screening = %+v", breakdown.Screening)
	}

	if total, _ := Score(job, profile, nil, 2026); total != 80 {
		t.Fatalf("tanpa screening total = %d, want 80", total)
	}
}

func TestScoreLevel(t *testing.T) {
	cases := []struct {
		level      string
		graduation *int
		experience int
		want       int
	}{
		{"internship", pointer(2027), 0, 20},
		{"internship", pointer(2025), 4, 20},
		{"internship", pointer(2024), 1, 10},
		{"internship", pointer(2020), 2, 0},
		{"internship", nil, 0, 10},
		{"entry", pointer(2020), 0, 20},
		{"entry", pointer(2020), 1, 20},
		{"entry", pointer(2020), 2, 10},
		{"entry", pointer(2026), 5, 10},
		{"entry", pointer(2020), 4, 0},
		{"junior", pointer(2020), 1, 20},
		{"junior", pointer(2020), 3, 20},
		{"junior", pointer(2020), 0, 10},
		{"junior", pointer(2020), 5, 10},
		{"junior", pointer(2020), 7, 0},
		{"mid", pointer(2020), 3, 20},
		{"mid", pointer(2020), 6, 20},
		{"mid", pointer(2020), 2, 10},
		{"mid", pointer(2020), 0, 0},
		{"mid", pointer(2020), 7, 0},
	}
	for _, test := range cases {
		profile := models.SeekerProfile{GraduationYear: test.graduation, ExperienceYears: test.experience}
		if got := scoreLevel(test.level, profile, 2026).Points; got != test.want {
			t.Errorf("%s, pengalaman %d: %d, want %d", test.level, test.experience, got, test.want)
		}
	}
}

func TestEvaluate(t *testing.T) {
	cases := []struct {
		questionType  string
		expected      *string
		answer        string
		passed, valid bool
	}{
		{"yes_no", pointer("yes"), "yes", true, true},
		{"yes_no", pointer("yes"), "no", false, true},
		{"yes_no", nil, "no", true, true},
		{"yes_no", pointer("yes"), "mungkin", false, false},
		{"number", pointer(">=2"), "2", true, true},
		{"number", pointer(">=2"), "1.5", false, true},
		{"number", pointer("<=5"), "5", true, true},
		{"number", pointer("<=5"), "6", false, true},
		{"number", nil, "6", true, true},
		{"number", pointer(">=2"), "banyak", false, false},
		{"number", pointer(">=2"), "Inf", false, false},
		{"number", nil, "NaN", false, false},
		{"text", nil, "halo", true, true},
		{"text", nil, "", false, true},
	}
	for _, test := range cases {
		question := models.ScreeningQuestion{Type: test.questionType, ExpectedAnswer: test.expected}
		passed, valid := Evaluate(question, test.answer)
		if passed != test.passed || valid != test.valid {
			t.Errorf("%s %q: (%v, %v), want (%v, %v)", test.questionType, test.answer, passed, valid, test.passed, test.valid)
		}
	}
}

func TestValidExpected(t *testing.T) {
	cases := []struct {
		questionType string
		expected     *string
		want         bool
	}{
		{"yes_no", pointer("yes"), true},
		{"yes_no", pointer(">=1"), false},
		{"number", pointer(">=1"), true},
		{"number", pointer(">1"), false},
		{"text", pointer("apa saja"), false},
		{"text", nil, true},
	}
	for _, test := range cases {
		if got := ValidExpected(test.questionType, test.expected); got != test.want {
			t.Errorf("%s: %v, want %v", test.questionType, got, test.want)
		}
	}
}
