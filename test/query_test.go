package objectbox_test

import (
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/test/assert"
	"github.com/objectbox/objectbox-go/test/model/iot"
	"github.com/objectbox/objectbox-go/test/model/iot/object"
	"testing"
)

func TestQueryBuilder(t *testing.T) {
	objectBox := iot.CreateObjectBox()
	box := objectBox.Box(1)
	box.RemoveAll()

	qb, err := objectBox.Query(1)
	assert.NoErr(t, err)
	query, err := qb.BuildAndDestroy()
	assert.NoErr(t, err)
	defer query.Destroy()

	objectBox.RunWithCursor(1, true, func(cursor *objectbox.Cursor) (err error) {
		bytesArray, err := query.FindBytes(cursor)
		assert.NoErr(t, err)
		assert.EqInt(t, 0, len(bytesArray.BytesArray))

		slice, err := query.Find(cursor)
		assert.NoErr(t, err)
		assert.EqInt(t, 0, len(slice.([]object.Event)))
		return
	})

	event := object.Event{
		Device: "dev1",
	}
	id1, err := box.Put(&event)
	assert.NoErr(t, err)

	event.Device = "dev2"
	id2, err := box.Put(&event)
	assert.NoErr(t, err)

	objectBox.RunWithCursor(1, true, func(cursor *objectbox.Cursor) (err error) {
		bytesArray, err := query.FindBytes(cursor)
		assert.NoErr(t, err)
		assert.EqInt(t, 2, len(bytesArray.BytesArray))

		slice, err := query.Find(cursor)
		assert.NoErr(t, err)
		events := slice.([]object.Event)
		if len(events) != 2 {
			t.Fatalf("unexpected size")
		}

		assert.Eq(t, id1, events[0].Id)
		assert.EqString(t, "dev1", events[0].Device)

		assert.Eq(t, id2, events[1].Id)
		assert.EqString(t, "dev2", events[1].Device)

		return
	})
}

func TestQueryBuilder_StringEq(t *testing.T) {
	objectBox := iot.CreateObjectBox()
	box := objectBox.Box(1)
	box.RemoveAll()

	iot.PutEvents(objectBox, 3)

	qb, err := objectBox.Query(1)
	assert.NoErr(t, err)
	defer qb.Destroy()
	qb.StringEq(2, "device 2", false)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Destroy()

	objectBox.RunWithCursor(1, true, func(cursor *objectbox.Cursor) (err error) {
		slice, err := query.Find(cursor)
		assert.NoErr(t, err)
		events := slice.([]object.Event)
		assert.EqInt(t, 1, len(events))
		assert.EqString(t, "device 2", events[0].Device)
		return
	})
}

func TestQueryBuilder_IntBetween(t *testing.T) {
	objectBox := iot.CreateObjectBox()
	box := objectBox.Box(1)
	box.RemoveAll()

	objectBox.SetDebugFlags(objectbox.DebugFlags_LOG_QUERIES | objectbox.DebugFlags_LOG_QUERY_PARAMETERS)

	events := iot.PutEvents(objectBox, 6)

	qb, err := objectBox.Query(1)
	assert.NoErr(t, err)
	defer qb.Destroy()
	start := events[2].Date
	end := events[4].Date
	qb.IntBetween(3, start, end)
	query, err := qb.Build()
	assert.NoErr(t, err)
	defer query.Destroy()

	objectBox.RunWithCursor(1, true, func(cursor *objectbox.Cursor) (err error) {
		slice, err := query.Find(cursor)
		assert.NoErr(t, err)
		events := slice.([]object.Event)
		assert.EqInt(t, 3, len(events))
		assert.Eq(t, start, events[0].Date)
		assert.Eq(t, start+1, events[1].Date)
		assert.Eq(t, end, events[2].Date)
		return
	})
}