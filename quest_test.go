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
	if quest.IsFinished() {
		t.Fatal("If last mission has been reached but not passed IsFinished() must return false")
	}

	quest.PassCurrent()
	if !quest.IsFinished() {
		t.Fatal("If last mission has been passed IsFinished() must return true")
	}
}

func TestEndHook(t *testing.T) {
	firstMissionEndHookCalled := false
	secondMissionEndHookCalled := false
	quest := NewQuest()
	quest.AddMission(Mission{
		Name: "1",
		Next: "2",
		End: func() {
			firstMissionEndHookCalled = true
		},
	}).AddMission(Mission{
		Name: "2",
		End: func() {
			secondMissionEndHookCalled = true
		},
	})

	quest.Start()
	quest.PassCurrent()

	if !firstMissionEndHookCalled {
		t.Fatal("Must call End() hook on first mission finish")
	}

	quest.PassCurrent()

	if !secondMissionEndHookCalled {
		t.Fatal("Must call End() hook on second mission finish")
	}
}

func TestForEachMission(t *testing.T) {
	count := 0
	quest := NewQuest()
	quest.AddMission(Mission{
		Name: "1",
		Next: "2",
	}).AddMission(Mission{
		Name: "2",
		Next: "3",
	}).AddMission(Mission{
		Name: "3",
		Next: "4",
	}).AddMission(Mission{
		Name: "4",
		Next: "5",
	}).AddMission(Mission{
		Name: "5",
		Next: "6",
	})

	quest.ForEachMission(func(m *Mission) {
		count++
	})

	if count != 5 {
		t.Fatal("ForEachMission() must iterate over each mission in quest")
	}
}

func TestResolveCurrentTask(t *testing.T) {
	q := NewQuest()
	q.AddMission(Mission{
		Tasks: Tasks{
			{
				Statement: "task 1",
			},
		},
	}).Start()

	_, err := q.ResolveCurrentTask("right answer")
	if err != nil {
		t.Fatal(err)
	}

	if !q.IsFinished() {
		t.Fatal("Quest mast be finished after the last task has been resolved")
	}
}
