package typing

// temp, dilemma

func TranslateStreamEncode(encodeCode string) (string) {
	switch encodeCode {
	case "1": return "H.264"
	case "2": return "H.265"
	}

	return ""
}

func TranslateResolution(ratioCode, resCode string) (w, h int) {
	switch ratioCode {
	case "4_3":
		switch resCode {
		case "320" : return 320, 240
		case "400" : return 400, 300
		case "640" : return 640, 480
		case "1280": return 1280, 960
		case "2048": return 2048, 1536
		case "800" : return 800, 600
		case "1600": return 1600, 1200
		case "2560": return 2560, 1920
		case "3072": return 3072, 2304
		}

	case "16_9":
		switch resCode {
		case "320" : return 320, 180
		case "640" : return 640, 360
		case "1280": return 1280, 720
		case "1920": return 1920, 1080
		case "2560": return 2560, 1440
		case "3072": return 3072, 1728
		case "3840": return 3840, 2160
		}

	case "1_1":
		switch resCode {
		case "320" : return 320, 320
		case "640" : return 640, 640
		case "1280": return 1280, 1280
		case "2192": return 2192, 2192
		case "2992": return 2992, 2992
		}
	}

	return 0, 0
}

func TranslateBitrate(bitrateCode string) (int) {
	kbpsMultiplier := 1000 // or 1024 ?

	switch bitrateCode {
	case "64"   : return 64 * kbpsMultiplier
	case "128"  : return 128 * kbpsMultiplier
	case "256"  : return 256 * kbpsMultiplier
	case "384"  : return 384 * kbpsMultiplier
	case "512"  : return 512 * kbpsMultiplier
	case "768"  : return 768 * kbpsMultiplier
	case "1024" : return 1024 * kbpsMultiplier
	case "1536" : return 1536 * kbpsMultiplier
	case "2048" : return 2048 * kbpsMultiplier
	case "3072" : return 3072 * kbpsMultiplier
	case "4096" : return 4096 * kbpsMultiplier
	case "6144" : return 6144 * kbpsMultiplier
	case "8192" : return 8192 * kbpsMultiplier
	case "10240": return 10240 * kbpsMultiplier
	case "12288": return 12288 * kbpsMultiplier
	case "14336": return 14336 * kbpsMultiplier
	case "16384": return 16384 * kbpsMultiplier
	case "20480": return 20480 * kbpsMultiplier
	case "24576": return 24576 * kbpsMultiplier
	}

	return 0
}

