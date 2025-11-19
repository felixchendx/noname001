package live

import "fmt"

// TODO: drawtext error in RHEL's ffmpeg (rhel9.6)

/*
	original credits: FernandoSiahaan1999

   The box below represents the text format for the stream frame:
   ________________________________
   | TLT          TCT          TRT |
   |                               |
   |                               |
   |                               |
   | TLM          TCM          TRM |
   |                               |
   |                               |
   |                               |
   | TLB          TCB          TRB |
   ---------------------------------

   Notes:
   - TLT : Text Left Top        (func prepTextLeftTop)
   - TLM : Text Left Middle     (func prepTextLeftMiddle)
   - TLB : Text Left Bottom     (func prepTextLeftBottom)

   - TCT : Text Center Top      (func prepTextCenterTop)
   - TCM : Text Center Middle   (func prepTextCenterMiddle)
   - TCB : Text Center Bottom   (func prepTextCenterBottom)

   - TRT : Text Right Top       (func prepTextRightTop)
   - TRM : Text Right Middle    (func prepTextRightMiddle)
   - TRB : Text Right Bottom    (func prepTextRightBottom)

   Description:
   This function sets the command for `drawtext` in FFMPEG based on the specified position within the frame.

   Parameters:
       - texts      []string   A collection of strings to be displayed in the frame (e.g., ["text1", "device code = dev001", ...])
       - fontColor  string      The color for the `drawtext` formatting (e.g., "black@0.9")
       - fontSize   string      The size of the font for `drawtext` (e.g., "11")

   Return:
       - string     The complete `drawtext` command string for FFMPEG (e.g.,
       "drawtext=text='[text1]':fontcolor=yellow@0.9:fontsize=0.02*h:x=0.01*W:y=0.98*h-0.02*h*2,drawtext=text='[device code = dev001]':fontcolor=yellow@0.9:fontsize=0.02*h:x=0.01*W:y=0.98*h-0.02*h*1")
*/

// =========================== DRAW TEXT LEFT POSITION - start =========================== //

func prepTextLeftTop(texts []string, fontColor string, fontSize string) string {
	var space string = fontSize
	defaultSetText := fmt.Sprintf(":fontcolor=%s:fontsize=%s:x=0.01*W:y=0.01*H+%s", fontColor, fontSize, space)

	return prepTextFormat(texts, fontColor, fontSize, defaultSetText, false)
}

func prepTextLeftMiddle(texts []string, fontColor string, fontSize string) string {
	amountText := len(texts)
	var space string = fontSize
	lineHeight := 30 
	startY := (30 * amountText) / 2 
	defaultSetText := fmt.Sprintf(":fontcolor=%s:fontsize=%s:x=0.01*W:y=(H/2-%d)+(%d)*%d+%s", fontColor, fontSize, startY, 0, lineHeight,space)

	return prepTextFormat(texts, fontColor, fontSize, defaultSetText, false)
}

func prepTextLeftBottom(texts []string, fontColor string, fontSize string) string {
	var space string = fontSize
	defaultSetText := fmt.Sprintf(":fontcolor=%s:fontsize=%s:x=0.01*W:y=0.98*h-%s", fontColor, fontSize, space)

	return prepTextFormat(texts, fontColor, fontSize, defaultSetText, true)
}
// =========================== DRAW TEXT LEFT POSITION - end =========================== //



// =========================== TEXT CENTER POSITION - start =========================== //

func prepTextCenterTop(texts []string, fontColor string, fontSize string) string {
	var space string = fontSize
	defaultSetText := fmt.Sprintf(":fontcolor=%s:fontsize=%s:x=W/2-text_w/2:y=0.01*H+%s", fontColor, fontSize, space)

	return prepTextFormat(texts, fontColor, fontSize, defaultSetText, false)
}

func prepTextCenterMiddle(texts []string, fontColor string, fontSize string) string {
	amountText := len(texts)
	var space string = fontSize
	lineHeight := 30 
	startY := (30 * amountText) / 2 
	defaultSetText := fmt.Sprintf(":fontcolor=%s:fontsize=%s:x=W/2-text_w/2:y=(H/2-%d)+(%d)*%d+%s", fontColor, fontSize, startY, 0, lineHeight,space)

	return prepTextFormat(texts, fontColor, fontSize, defaultSetText, false)
}

func prepTextCenterBottom(texts []string, fontColor string, fontSize string) string {
	var space string = fontSize
	defaultSetText := fmt.Sprintf(":fontcolor=%s:fontsize=%s:x=W/2-text_w/2:y=0.98*H-%s", fontColor, fontSize, space)

	return prepTextFormat(texts, fontColor, fontSize, defaultSetText, true)
}
// =========================== DRAW TEXT CENTER POSITION - end =========================== //



// =========================== DRAW TEXT RIGHT POSITION - start =========================== //

func prepTextRightTop(texts []string, fontColor string, fontSize string) string {
	var space string = fontSize
	defaultSetText := fmt.Sprintf(":fontcolor=%s:fontsize=%s:x=W-text_w-10:y=0.01*H+%s", fontColor, fontSize, space)

	return prepTextFormat(texts, fontColor, fontSize, defaultSetText, false)
}

func prepTextRightMiddle(texts []string, fontColor string, fontSize string) string {
	var space string = fontSize
	amountText := len(texts)
	lineHeight := 30 
	startY := (30 * amountText) / 2 
	defaultSetText := fmt.Sprintf(":fontcolor=%s:fontsize=%s:x=W-text_w-10:y=(H/2-%d)+(%d)*%d+%s", fontColor, fontSize, startY, 0, lineHeight,space)
	
	return prepTextFormat(texts, fontColor, fontSize, defaultSetText, false)
}


func prepTextRightBottom(texts []string, fontColor string, fontSize string) string {
	var space string = fontSize
	defaultSetText := fmt.Sprintf(":fontcolor=%s:fontsize=%s:x=W-text_w-10:y=0.98*H-%s", fontColor, fontSize, space)
	
	return prepTextFormat(texts, fontColor, fontSize, defaultSetText, true)
}
// =========================== DRAW TEXT RIGHT POSITION - end =========================== //


func prepTextFormat(texts []string, fontColor string, fontSize string, formatText string, flagBottom bool) string {
	amountText := len(texts)
	if amountText <= 0 {
		return ""
	}

	var cmdText string = ""
	for idx := 0; idx < amountText; idx++ {
		if flagBottom{
			cmdText += fmt.Sprintf(`drawtext=text='%s'%s*%d`, texts[idx], formatText, (amountText-idx))
		} else {
			cmdText += fmt.Sprintf(`drawtext=text='%s'%s*%d`, texts[idx], formatText, idx)
		}

		if idx < amountText-1 {
			cmdText += ","
		}
	}
	return cmdText
}
