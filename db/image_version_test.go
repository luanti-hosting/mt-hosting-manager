package db_test

import (
	"mt-hosting-manager/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageVersionRepository(t *testing.T) {
	repos := SetupRepos(t)

	list, err := repos.ImageVersionRepo.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.True(t, len(list) > 0)

	mtui_version, err := repos.ImageVersionRepo.GetByName(types.ImageNameUI)
	assert.NoError(t, err)
	assert.NotNil(t, mtui_version)
	assert.NotEqual(t, "", mtui_version.Name)
	assert.NotEqual(t, "", mtui_version.Version)
}
