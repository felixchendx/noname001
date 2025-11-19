'use strict';

// TODO: sweep
const nsMain = {};
nsMain.focusOn = (_elID) => { const el = document.getElementById(_elID); if (el != null) { el.focus(); } };

const _nsMain = {};
// temp, make date util stuffs
// https://moment.github.io/luxon/#/formatting?id=table-of-tokens
_nsMain.formatTimestamp01 = (ts) => { return luxon.DateTime.fromISO(ts).toFormat('dd LLLL yyyy, HH:mm:ss ZZZZ'); };
