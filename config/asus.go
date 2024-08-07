package config

var (
	ASUS_CODES = map[string]string{
		"NB": "NB",
		"NR": "NR",
		"NX": "NX",
		"PF": "PF",
		"PT": "PT",
	}
	ASUS_LONGNAMES = map[string]string{
		"LWEP":          "Garantía extendida",
		"ADP":           "Protección contra daño accidental",
		"BSP":           "Protección para batería",
		"HDD_RETENTION": "Protección de disco duro",
		"OSS":           "Servicio domicilio",
	}
	ASUS_SHORTNAMES = map[string]string{
		"LWEP":          "Garantía",
		"ADP":           "Daño accidental",
		"BSP":           "Batería",
		"HDD_RETENTION": "Disco duro",
		"OSS":           "Domicilio",
	}
)
