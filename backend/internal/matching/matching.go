package matching

import (
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"karsa/internal/models"
)

const (
	skillPoints     = 60
	levelPoints     = 20
	screeningPoints = 20
)

var (
	levels          = []string{"internship", "entry", "junior", "mid"}
	numberCondition = regexp.MustCompile(`^(>=|<=)(\d{1,9})$`)
)

func Score(job models.Job, profile models.SeekerProfile, passed []bool, currentYear int) (int, models.ScoreBreakdown) {
	breakdown := models.ScoreBreakdown{
		Skill:     scoreSkills(job.Skills, profile.Skills),
		Level:     scoreLevel(job.Level, profile, currentYear),
		Screening: scoreScreening(passed),
	}
	return breakdown.Skill.Points + breakdown.Level.Points + breakdown.Screening.Points, breakdown
}

func scoreSkills(needed, owned []string) models.SkillScore {
	result := models.SkillScore{Matched: []string{}, Missing: []string{}}
	for _, skill := range needed {
		if slices.ContainsFunc(owned, func(own string) bool { return strings.EqualFold(own, skill) }) {
			result.Matched = append(result.Matched, skill)
		} else {
			result.Missing = append(result.Missing, skill)
		}
	}
	result.Points = portion(skillPoints, len(result.Matched), len(needed))
	return result
}

func scoreLevel(level string, profile models.SeekerProfile, currentYear int) models.LevelScore {
	if fitsLevel(level, profile, currentYear) {
		return models.LevelScore{Points: levelPoints, Note: "Sesuai level"}
	}
	position := slices.Index(levels, level)
	for index, neighbor := range levels {
		if (index == position-1 || index == position+1) && fitsLevel(neighbor, profile, currentYear) {
			return models.LevelScore{Points: levelPoints / 2, Note: "Selisih satu tingkat"}
		}
	}
	return models.LevelScore{Note: "Di luar rentang level"}
}

func fitsLevel(level string, profile models.SeekerProfile, currentYear int) bool {
	years := profile.ExperienceYears
	switch level {
	case "internship":
		return profile.GraduationYear != nil && *profile.GraduationYear >= currentYear-1
	case "entry":
		return years <= 1
	case "junior":
		return years >= 1 && years <= 3
	case "mid":
		return years >= 3 && years <= 6
	}
	return false
}

func scoreScreening(passed []bool) models.ScreeningScore {
	result := models.ScreeningScore{Total: len(passed)}
	for _, isPassed := range passed {
		if isPassed {
			result.Passed++
		}
	}
	result.Points = portion(screeningPoints, result.Passed, result.Total)
	return result
}

func portion(points, part, total int) int {
	if total == 0 {
		return points
	}
	return int(math.Round(float64(points*part) / float64(total)))
}

func ValidExpected(questionType string, expected *string) bool {
	if expected == nil {
		return true
	}
	switch questionType {
	case "yes_no":
		return *expected == "yes" || *expected == "no"
	case "number":
		return numberCondition.MatchString(*expected)
	}
	return false
}

func Evaluate(question models.ScreeningQuestion, answer string) (passed, valid bool) {
	if answer == "" {
		return false, true
	}
	switch question.Type {
	case "yes_no":
		valid = answer == "yes" || answer == "no"
		return valid && (question.ExpectedAnswer == nil || answer == *question.ExpectedAnswer), valid
	case "number":
		number, err := strconv.ParseFloat(answer, 64)
		if err != nil || math.IsNaN(number) || math.IsInf(number, 0) {
			return false, false
		}
		return meetsCondition(number, question.ExpectedAnswer), true
	}
	return true, true
}

func meetsCondition(number float64, expected *string) bool {
	if expected == nil {
		return true
	}
	parts := numberCondition.FindStringSubmatch(*expected)
	if parts == nil {
		return false
	}
	limit, _ := strconv.ParseFloat(parts[2], 64)
	if parts[1] == ">=" {
		return number >= limit
	}
	return number <= limit
}
