package constant

type Brand string
const (
	BRAND__DAHUA     Brand = "dahua"
	BRAND__HIKVISION Brand = "hikvision"
	BRAND__HUAWEI    Brand = "huawei"

	// temp
	BRAND__PANASONIC_NETCAM Brand = "panasonic-netcam"
)

var (
	BRANDS = []Brand{
		BRAND__DAHUA,
		BRAND__HIKVISION,
		BRAND__HUAWEI,

		BRAND__PANASONIC_NETCAM,
	}
)

// temp
type BrandStreamType string

const (
	BRAND_STREAM_TYPE__NONE BrandStreamType = "none"

	DAHUA__MAIN_STREAM    BrandStreamType = "main"
	DAHUA__EXTRA_STREAM_1 BrandStreamType = "extra1"
	DAHUA__EXTRA_STREAM_2 BrandStreamType = "extra2"
	DAHUA__EXTRA_STREAM_3 BrandStreamType = "extra3"

	HIKVISION__MAIN_STREAM  BrandStreamType = "main"
	HIKVISION__SUB_STREAM   BrandStreamType = "sub"
	HIKVISION__THIRD_STREAM BrandStreamType = "third"

	PANASONIC_NETCAM__STREAM_1 BrandStreamType = "stream1"
	PANASONIC_NETCAM__STREAM_2 BrandStreamType = "stream2"
	PANASONIC_NETCAM__STREAM_3 BrandStreamType = "stream3"
	PANASONIC_NETCAM__STREAM_4 BrandStreamType = "stream4"
)

var (
	DAHUA__STREAM_TYPES = map[BrandStreamType]string{
		DAHUA__MAIN_STREAM   : "Main Stream",
		DAHUA__EXTRA_STREAM_1: "Extra Stream 1",
		DAHUA__EXTRA_STREAM_2: "Extra Stream 2",
		DAHUA__EXTRA_STREAM_3: "Extra Stream 3",
	}

	HIKVISION__STREAM_TYPES = map[BrandStreamType]string{
		HIKVISION__MAIN_STREAM : "Main Stream",
		HIKVISION__SUB_STREAM  : "Sub Stream",
		HIKVISION__THIRD_STREAM: "Third Stream",
	}

	PANASONIC_NETCAM__STREAM_TYPES = map[BrandStreamType]string{
		PANASONIC_NETCAM__STREAM_1: "Stream(1)",
		PANASONIC_NETCAM__STREAM_2: "Stream(2)",
		PANASONIC_NETCAM__STREAM_3: "Stream(3)",
		PANASONIC_NETCAM__STREAM_4: "Stream(4)",
	}
)

// func BrandStreamTypeFromString(s string) (BrandStreamType) {
// 	switch s {
// 	}
// }
