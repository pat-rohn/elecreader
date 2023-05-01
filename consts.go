package elecreader

/*
F.F(00)                    // Error code (always first in the list)
0.0(            1844)		// Customer identification (16 char. string) (Readout)
C.1.0(14593907)				// Meter identification (8 character string)
C.1.1(        )				// Manufacture identification (8 char. string)
1.8.1(060740.806*kWh)		// Active energy import rate 1
1.8.2(043092.307*kWh)		// Active energy import rate 2
2.8.1(000000.000*kWh)		// Active energy export rate 1
2.8.2(000000.000*kWh)		// Active energy export rate 2
1.8.0(103833.113*kWh)		// Total active energy import
2.8.0(000000.000*kWh)		// Total active energy export
15.8.0(103833.114*kWh)		// Active Energy - Absolute Value [|+A|+|-A|] {+kWh}
C.7.0(002)					// Power off counter
32.7(233*V)					// Line voltage L1
52.7(231*V)
72.7(232*V)
31.7(000.78*A)				// Line current L1
51.7(000.49*A)
71.7(001.44*A)
82.8.1(0000)				// Terminal cover removal counter
82.8.2(0000)				// DC field detection counter
0.2.0(M26)					// Software version
C.5.0(0420)					// Status code (see section 5.5.3)
!<ETX><BCC>					// End of text, Checksum
*/

type Result struct {
	TotalActiveEnergyImport float64 // kWh
	ActiveEnergyImportRate1 float64 // kWh
	ActiveEnergyImportRate2 float64 // kWh
	TotalActiveEnergyExport float64 // kWh
	ActiveEnergyExportRate1 float64 // kWh
	ActiveEnergyExportRate2 float64 // kWh
	ActiveEnergyAbsolute    float64 // kWh
	CurrentLine1            float64 // A
	CurrentLine2            float64 // A
	CurrentLine3            float64 // A
	VoltageLine1            float64 // V
	VoltageLine2            float64 // V
	VoltageLine3            float64 // V
}
