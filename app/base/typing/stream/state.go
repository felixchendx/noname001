package stream

type (
	LiveState         string
	LiveFailState     string
	LiveStreamerState string
)

const (
	LIVE_STATE__NEW          LiveState = "ls:new"
	LIVE_STATE__INACTIVE     LiveState = "ls:inactive"
	LIVE_STATE__INIT_BEGIN   LiveState = "ls:init:begin"
	LIVE_STATE__INIT_FAIL    LiveState = "ls:init:fail"
	LIVE_STATE__INIT_OK      LiveState = "ls:init:ok"
	LIVE_STATE__RELOAD_BEGIN LiveState = "ls:reload:begin"
	LIVE_STATE__RELOAD_FAIL  LiveState = "ls:reload:fail"
	LIVE_STATE__RELOAD_OK    LiveState = "ls:reload:ok"
	LIVE_STATE__DESTROY      LiveState = "ls:destroy"
	// umbrella fail state for bg stuffs
	LIVE_STATE__BG_FAIL      LiveState = "ls:bg:fail"

	LIVE_FAIL_STATE__NONE        LiveFailState = "lfs:none"
	LIVE_FAIL_STATE__DB          LiveFailState = "lfs:db"
	LIVE_FAIL_STATE__DEP_PROFILE LiveFailState = "lfs:dep:profile"
	LIVE_FAIL_STATE__DEP_DEVICE  LiveFailState = "lfs:dep:device"
	LIVE_FAIL_STATE__DEP_FILE    LiveFailState = "lfs:dep:file"
	LIVE_FAIL_STATE__OTHER       LiveFailState = "lfs:other"
	LIVE_FAIL_STATE__BG_STREAMER LiveFailState = "lfs:bg:streamer"

	LIVE_STREAMER_STATE__NEW             LiveStreamerState = "lss:new"
	LIVE_STREAMER_STATE__START           LiveStreamerState = "lss:start"
	LIVE_STREAMER_STATE__RUNNING         LiveStreamerState = "lss:running"
	LIVE_STREAMER_STATE__STOP_NORMAL     LiveStreamerState = "lss:stop:normal"
	LIVE_STREAMER_STATE__STOP_UNEXPECTED LiveStreamerState = "lss:stop:unexpected"
)
