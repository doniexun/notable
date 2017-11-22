package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecondaryCreateNew(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	note := Note{
		Subject: fmt.Sprintf("secondary note creation %s", time.Now()),
	}
	_, err := mock.secondary.update(note)
	assert.Nil(t, err)
}

func TestSecondaryList(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	note := Note{
		Subject: fmt.Sprintf("secondary note creation %s", time.Now()),
	}
	_, err := mock.secondary.update(note)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(mock.secondary.list()))
	assert.Equal(t, note.Subject, mock.secondary.list()[0].Subject)
}

func TestSecondaryUpdateExisting(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	_, s, _, _ := createTestNote(mock, "")
	s.Content = fmt.Sprintf("secondary note mutation %s", time.Now())
	_, err := mock.secondary.update(s)
	assert.Nil(t, err)
}

func TestSecondaryReadAfterWrite(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	_, s, _, _ := createTestNote(mock, "")
	s.Content = fmt.Sprintf("secondary read after write %s", time.Now())
	_, err := mock.secondary.update(s)
	assert.Nil(t, err)
	note, err := mock.secondary.getNoteByUID(s.UID, "")
	assert.Nil(t, err)
	assert.Equal(t, s.Content, note.Content)
}

func TestPrimaryReadAfterSecondaryWrite(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	_, s, _, _ := createTestNote(mock, "")
	s.Content = fmt.Sprintf("secondary note mutation for listing %s", time.Now())
	_, err := mock.secondary.update(s)
	assert.Nil(t, err)

	// Simulare the recover process the primary would be running
	consumeSecondaries(mock.db, Secondary{Path: mock.db.dbFilePath()}, new(messenger))

	// Read through the primmary to see if it finds the changes
	p, err := db.getNoteByUID(s.UID, "")
	assert.Nil(t, err)
	assert.Equal(t, s.Content, p.Content)
}

func TestReloadAsNeeded(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	_, note, _, _ := createTestNote(mock, "")
	frontend, backend := new(messenger), new(messenger)
	frontendCh := frontend.add()
	backendCh := backend.add()
	go reloadAsNeeded(mock.db, frontend, backend)
	time.Sleep(time.Second * 3)
	mock.db.update(note)
	assert.Equal(t, "reload", <-frontendCh)
	backend.send("stop")
	<-backendCh
}

func TestSecondaryUpdate(t *testing.T) {
	mock := setup(t)
	defer tearDown(mock)
	note := Note{
		Subject: "subject",
		Content: "body",
	}
	persisted, err := mock.secondary.update(note)
	assert.Nil(t, err)
	assert.Equal(t, note.Subject, persisted.Subject)
}