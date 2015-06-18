package msg

import (
	"log"
)

type VolcanicAlert struct {
	Level    int
	Activity string
	Hazards  string
}

var VolcanicAlertLevels = [...]VolcanicAlert{
	VolcanicAlert{
		Level:    0,
		Activity: `No volcanic unrest.`,
		Hazards:  `Volcanic environment hazards.`,
	},
	VolcanicAlert{
		Level:    1,
		Activity: `Minor volcanic unrest.`,
		Hazards:  `Volcanic unrest hazards.`,
	},
	VolcanicAlert{
		Level:    2,
		Activity: `Moderate to heightened volcanic unrest.`,
		Hazards:  `Volcanic unrest hazards, potential for eruption hazards.`,
	},
	VolcanicAlert{
		Level:    3,
		Activity: `Minor volcanic eruption.`,
		Hazards:  `Eruption hazards near vent. Note: ash, lava flow, and lahar (mudflow) hazards may impact areas distant from the volcano.`,
	},
	VolcanicAlert{
		Level:    4,
		Activity: `Moderate volcanic eruption.`,
		Hazards:  `Eruption hazards on and near volcano. Note: ash, lava flow, and lahar (mudflow) hazards may impact areas distant from the volcano.`,
	},
	VolcanicAlert{
		Level:    5,
		Activity: `Major volcanic eruption.`,
		Hazards:  `Eruption hazards on and beyond volcano. Note: ash, lava flow, and lahar (mudflow) hazards may impact areas distant from the volcano.`,
	},
}

type Volcano struct {
	VolcanoID string
	Title     string
}

var Volcanoes = [...]Volcano{
	Volcano{
		VolcanoID: `aucklandvolcanicfield`,
		Title:     `Auckland Volcanic Field`,
	},
	Volcano{
		VolcanoID: `kermadecislands`,
		Title:     `Kermadec Islands`,
	},
	Volcano{
		VolcanoID: `mayorisland`,
		Title:     `Mayor Island`,
	},
	Volcano{
		VolcanoID: `ngauruhoe`,
		Title:     `Ngauruhoe`,
	},
	Volcano{
		VolcanoID: `northland`,
		Title:     `Northland`,
	},
	Volcano{
		VolcanoID: `okataina`,
		Title:     `Okataina`,
	},
	Volcano{
		VolcanoID: `rotorua`,
		Title:     `Rotorua`,
	},
	Volcano{
		VolcanoID: `ruapehu`,
		Title:     `Ruapehu`,
	},
	Volcano{
		VolcanoID: `taupo`,
		Title:     `Taupo`,
	},
	Volcano{
		VolcanoID: `tongariro`,
		Title:     `Tongariro`,
	},
	Volcano{
		VolcanoID: `taranakiegmont`,
		Title:     `Taranaki/Egmont`,
	},
	Volcano{
		VolcanoID: `whiteisland`,
		Title:     `White Island`,
	},
}

/*
For sending Volcanic Alert Level messages.
See http://info.geonet.org.nz/x/PYAO
*/
type VAL struct {
	Volcano       Volcano
	VolcanicAlert VolcanicAlert
	err           error
}

func (v *VAL) Err() error {
	return v.err
}

func (v *VAL) SetErr(err error) {
	v.err = err
}

func (v *VAL) RxLog() {
	if v.err != nil {
		return
	}

	log.Printf("Received volcanic alert level update: %s %d", v.Volcano.Title, v.VolcanicAlert.Level)
}

func (v *VAL) TxLog() {
	if v.err != nil {
		return
	}

	log.Printf("Sending volcanic alert level update: %s %d", v.Volcano.Title, v.VolcanicAlert.Level)
}
