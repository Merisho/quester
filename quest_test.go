package quester

import (
	"testing"
)

func TestAddMission(t *testing.T) {
	quest := NewQuest()
	quest.AddMission(Mission{})

	if quest.Length() != 1 {
		t.Fatal("Quest length must be equal 1 after mission was added")
	}
}

func TestNoMissionsErr(t *testing.T) {
	quest := NewQuest()
	err := quest.Start()

	if err != NoMissionsErr {
		t.Fatal("Quest.Start() must return NoMissionsErr in case quest started without missions")
	}
}

func TestCurrentMission(t *testing.T) {
	quest := NewQuest()
	quest.AddMission(Mission{
		Name: "mission 1",
		Next: "mission 2",
	})
	quest.AddMission(Mission{
		Name: "mission 2",
	})
	quest.Start()

	if quest.Current().Name != "mission 1" {
		t.Fatal("First added mission must be the current")
	}
}

func TestStartMissionHook(t *testing.T) {
	missionStarted := false
	quest := NewQuest()
	quest.AddMission(Mission{
		Name: "mission 1",
		Start: func() {
			missionStarted = true
		},
	})

	quest.Start()

	if !missionStarted {
		t.Fatal("Mission.Start hook must be called on mission start")
	}
}

func TestNoMissionsStartError(t *testing.T) {
	quest := NewQuest()

	err := quest.Start()

	if err != NoMissionsErr {
		t.Fatal("Return value must be NoMissionErr if no missions are specified")
	}
}

func TestPassMission(t *testing.T) {
	nextMissStarted := false
	quest := NewQuest()
	quest.AddMission(Mission{
		Name: "mission 1",
		Next: "mission 2",
	})
	quest.AddMission(Mission{
		Name: "mission 2",
		Start: func() {
			nextMissStarted = true
		},
	})

	quest.Start()
	quest.PassCurrent()

	if quest.Current().Name != "mission 2" {
		t.Fatal("After current mission pass next mission must become current")
	}

	if !nextMissStarted {
		t.Fatal("Must call Start hook of next mission in case of current mission pass")
	}
}

func TestFinishedQuest(t *testing.T) {
	quest := NewQuest()
	quest.AddMission(Mission{
		Name: "mission 1",
		Next: "mission 2",
	}).AddMission(Mission{
		Name: "mission 2",
	})

	quest.Start()
	quest.PassCurrent()
	finished := quest.IsFinished()

	if !finished {
		t.Fatal("If last mission has been passed IsFinished() must return true")
	}
}

func TestEndHook(t *testing.T) {
	endHookCalled := false
	quest := NewQuest()
	quest.AddMission(Mission{
		Name: "1",
		Next: "2",
		End: func() {
			endHookCalled = true
		},
	}).AddMission(Mission{
		Name: "2",
	})

	quest.Start()
	quest.PassCurrent()

	if !endHookCalled {
		t.Fatal("Must call End() hook on mission leave")
	}
}

func TestQuestFlow(t *testing.T) {
	finalStarted := false
	initEnded := false
	quest := NewQuest()
	quest.AddMission(Mission{
		Name: "Init",
		Next: "Mission 2",
		End: func() {
			initEnded = true
		},
	}).AddMission(Mission{
		Name: "Mission 2",
		Next: "Mission 3",
	}).AddMission(Mission{
		Name: "Mission 3",
		Next: "Final",
	}).AddMission(Mission{
		Name: "Final",
		Start: func() {
			finalStarted = true
		},
	})

	err := quest.Start()
	if err != nil {
		t.Fatal(err)
	}
	if quest.Current().Name != "Init" {
		t.Fatal("First mission must be Init")
	}

	err = quest.PassCurrent()
	if err != nil {
		t.Fatal(err)
	}
	if quest.Current().Name != "Mission 2" {
		t.Fatal("Second mission must be Mission 2")
	}
	if !initEnded {
		t.Fatal("Must call End() hook on Init mission end")
	}

	err = quest.PassCurrent()
	if err != nil {
		t.Fatal(err)
	}
	if quest.Current().Name != "Mission 3" {
		t.Fatal("Third mission must be Mission 3")
	}

	err = quest.PassCurrent()
	if err != nil {
		t.Fatal(err)
	}
	if quest.Current().Name != "Final" {
		t.Fatal("Last mission must be Final")
	}
	if !finalStarted {
		t.Fatal("Start hook on final mission must be called")
	}
}
