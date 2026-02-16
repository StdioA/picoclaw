package skills

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkillsInfoValidate(t *testing.T) {
	testcases := []struct {
		name        string
		skillName   string
		description string
		wantErr     bool
		errContains []string
	}{
		{
			name:        "valid-skill",
			skillName:   "valid-skill",
			description: "a valid skill description",
			wantErr:     false,
		},
		{
			name:        "empty-name",
			skillName:   "",
			description: "description without name",
			wantErr:     true,
			errContains: []string{"name is required"},
		},
		{
			name:        "empty-description",
			skillName:   "skill-without-description",
			description: "",
			wantErr:     true,
			errContains: []string{"description is required"},
		},
		{
			name:        "empty-both",
			skillName:   "",
			description: "",
			wantErr:     true,
			errContains: []string{"name is required", "description is required"},
		},
		{
			name:        "name-with-spaces",
			skillName:   "skill with spaces",
			description: "invalid name with spaces",
			wantErr:     true,
			errContains: []string{"name must be alphanumeric with hyphens"},
		},
		{
			name:        "name-with-underscore",
			skillName:   "skill_underscore",
			description: "invalid name with underscore",
			wantErr:     true,
			errContains: []string{"name must be alphanumeric with hyphens"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			info := SkillInfo{
				Name:        tc.skillName,
				Description: tc.description,
			}
			err := info.validate()
			if tc.wantErr {
				assert.Error(t, err)
				for _, msg := range tc.errContains {
					assert.ErrorContains(t, err, msg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoadSkills(t *testing.T) {
	// Test loading weather skill from workspace skills directory
	workspace := "../../workspace"
	loader := NewSkillsLoader(workspace, "", "")

	// Test LoadSkill directly to verify content loading
	content, found := loader.LoadSkill("weather")
	assert.True(t, found, "weather skill content should be found")
	assert.NotEmpty(t, content, "weather skill content should not be empty")
	assert.Contains(t, content, "Weather", "content should contain 'Weather' heading")
	assert.Contains(t, content, "wttr.in", "content should mention wttr.in service")
	assert.Contains(t, content, "Open-Meteo", "content should mention Open-Meteo service")

	// Test that the SKILL.md file exists and can be read
	weatherSkillPath := "../../workspace/skills/weather/SKILL.md"
	skillData, err := os.ReadFile(weatherSkillPath)
	assert.NoError(t, err, "should be able to read weather SKILL.md file")
	assert.NotEmpty(t, skillData, "weather SKILL.md should not be empty")

	// Test ListSkills to verify weather skill is properly loaded with metadata
	skills := loader.ListSkills()
	var weatherSkill *SkillInfo
	for i := range skills {
		if skills[i].Name == "weather" {
			weatherSkill = &skills[i]
			break
		}
	}
	assert.NotNil(t, weatherSkill, "weather skill should be found in ListSkills")
	assert.Equal(t, "weather", weatherSkill.Name)
	assert.Equal(t, "Get current weather and forecasts (no API key required).", weatherSkill.Description)
	assert.Equal(t, "workspace", weatherSkill.Source)
	assert.Contains(t, weatherSkill.Path, "weather/SKILL.md")
}
