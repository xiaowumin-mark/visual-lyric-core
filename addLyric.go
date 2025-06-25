package main

import (
	"fmt"
	"log"
	"math"

	//"math"

	"syscall/js"
	"time"

	lyrics "github.com/xiaowumin-mark/visual-lyric-core/lyric"
)

func addLyric(index int, lrcs *lyrics.Lyrics) {

	addIndex(index)
	// 在js控制台打印nowPlayingIndex
	currentTime := getCurrentTime(audio)
	words := lrcs.Contents[index].Primary.Blocks

	//lineTimeE := math.Min(float64(lrcs.Contents[index].Primary.End), float64(lrcs.Contents[index].Primary.Blocks[len(lrcs.Contents[index].Primary.Blocks)-1].End))
	//lineTimeB := math.Max(float64(lrcs.Contents[index].Primary.Begin), float64(lrcs.Contents[index].Primary.Blocks[0].Begin))
	//lineTime := lineTimeE - lineTimeB
	lineTime := lrcs.Contents[index].Primary.End - lrcs.Contents[index].Primary.Begin
	fmt.Println(time.Duration(lineTime).Milliseconds())

	if len(lrcs.Contents[index].Backgrounds) != 0 {
		lrcs.Contents[index].BackgroundsEle.Get("style").Set("display", "block")
		GsetTimeout(func() {
			lrcs.Contents[index].BackgroundsEle.Get("classList").Call("add", "bgShow")

			GsetTimeout(func() {
				gd(bubbleSort(nowPlayingIndex)[0], lrcs, false)
			}, 50*time.Millisecond)
		}, 10*time.Millisecond)

		for _, item := range lrcs.Contents[index].Backgrounds {
			for _, word := range item.Blocks {
				intervalTime := word.Begin - currentTime
				//word.Ele.Get("style").Call("setProperty", "--p", "100%")
				//word.Ele.Get("style").Call("setProperty", "--rp", "140%")
				duration := word.End - word.Begin

				eleWidth := word.Ele.Get("offsetWidth").Float()

				oPeo := float64(bglw) + eleWidth

				V := eleWidth / duration.Seconds() // 原速度（100% / duration）

				// 计算 20px 过渡的额外时间
				//extraTime := (float64(bglw) * 1 / eleWidth * 100) / (oPeo / duration.Seconds())

				// 调整后的动画时间（保持速度 V 不变）
				nT := oPeo / V

				animation := gsap.Call("to", word.Ele, map[string]interface{}{
					// duration单位为秒
					//"duration": duration.Seconds() * 1.05,
					"duration": nT,
					"ease":     "none",
					//"--p":      "100%",
					//"--rp":     "140%",
					// 延时触发
					//"delay": intervalTime.Seconds() * 0.95,
					"--p": fmt.Sprintf("%vpx", oPeo),

					"delay": intervalTime.Seconds(),
				})
				word.GsapAnimation = animation

			}
		}
	} else {
		//if index != noePlayingOne {
		//	gd(bubbleSort(nowPlayingIndex)[0], lrcs, false)
		//	noePlayingOne = index
		//}
		//gd(bubbleSort(nowPlayingIndex)[0], lrcs, false)
		gd(bubbleSort(nowPlayingIndex)[0], lrcs, false)
		//hasScrolledInRemove = false // 重置标志位
	}
	var wordsN []*lyrics.Block

	for _, word := range words {
		if word.Begin == 0 && word.End == 0 {
			continue
		}
		wordsN = append(wordsN, word)
	}

	//curX := 0.0 // 累积偏移量（整行动画进度）
	for wi, word := range wordsN {
		//laspv := 0.0
		//for _, word := range words {
		//	laspv += math.Max(float64(word.End), laspv)
		//}
		//
		//totalFadeDuration := float64(lrcs.Contents[index].Primary.End) - float64(lrcs.Contents[index].Primary.Begin)
		////for _, w := range words {
		////	totalFadeDuration = math.Max(float64(w.End), totalFadeDuration)
		////}
		////totalFadeDuration = math.Max(float64(lrcs.Contents[index].Primary.End), totalFadeDuration)
		////totalFadeDuration -= float64(lrcs.Contents[index].Primary.Begin)
		////totalFadeDuration := math.Max(laspv, float64(lrcs.Contents[index].Primary.End)) - float64(lrcs.Contents[index].Primary.Begin)
		//wordEle := word.Ele
		//wordWidth := wordEle.Get("offsetWidth").Float()
		////wordHeight := wordEle.Get("offsetHeight").Float()
		//fadeWidth := wordWidth * fadeRatio
		//
		//widthBeforeSelf := 0.0
		//for i := 0; i < wi; i++ {
		//	widthBeforeSelf += words[i].Ele.Get("offsetWidth").Float()
		//}
		//minOffset := -(wordWidth + fadeWidth)
		//clampOffset := func(x float64) float64 {
		//	return math.Max(math.Min(x, 0), minOffset)
		//}
		//curPos := -widthBeforeSelf - wordWidth - fadeWidth
		//timeOffset := time.Duration(0)
		//
		//var frames []interface{}
		//
		//lastPos := curPos
		//lastTime := time.Duration(0)
		//
		//pushFrame := func() {
		//	moveOffset := curPos - lastPos
		//	time_ := math.Max(0, math.Min(1, float64(timeOffset)))
		//	duration := time_ - float64(lastTime)
		//	d := math.Abs(duration / moveOffset)
		//	if curPos > minOffset && lastPos < minOffset {
		//		staticTime := math.Abs(lastPos-minOffset) * d
		//		//value := `${clampOffset(lastPos)}px 0`;
		//		value := fmt.Sprintf("%vpx 0", clampOffset(curPos))
		//		//const frame: Keyframe = {
		//		//	offset: lastTime + staticTime,
		//		//	maskPosition: value,
		//		//};
		//		//frames.push(frame);
		//		frame := map[string]interface{}{
		//			"offset":             float64(lastTime) + staticTime,
		//			"backgroundPosition": value,
		//		}
		//
		//		frames = append(frames, frame)
		//	}
		//	if curPos > 0 && lastPos < 0 {
		//		staticTime := math.Abs(lastPos) * d
		//		value := fmt.Sprintf("%vpx 0", clampOffset(curPos))
		//		frame := map[string]interface{}{
		//			"offset":             float64(lastTime) + staticTime,
		//			"backgroundPosition": value,
		//		}
		//
		//		frames = append(frames, frame)
		//	}
		//
		//	//value := `${clampOffset(curPos)}px 0`
		//	value := fmt.Sprintf("%vpx 0", clampOffset(curPos))
		//	frame := map[string]interface{}{
		//		"offset":             time_,
		//		"backgroundPosition": value,
		//	}
		//
		//	frames = append(frames, frame)
		//	lastPos = curPos
		//	lastTime = time.Duration(time_)
		//
		//}
		//
		//pushFrame()
		//
		//lastTimeStamp := time.Duration(0)
		//
		//for j, otherWord := range words {
		//	curTimeStamp := otherWord.Begin - lrcs.Contents[index].Primary.Begin
		//	staticDuration := curTimeStamp - lastTimeStamp
		//	timeOffset += staticDuration / time.Duration(totalFadeDuration)
		//	if staticDuration > 0 {
		//		pushFrame()
		//	}
		//	lastTimeStamp = curTimeStamp
		//
		//	fadeDuration := otherWord.End - otherWord.Begin
		//	timeOffset += fadeDuration / time.Duration(totalFadeDuration)
		//	curPos += otherWord.Ele.Get("offsetWidth").Float()
		//	if j == 0 {
		//		curPos += fadeWidth * 1.5
		//	}
		//	if j == len(lrcs.Contents[index].Primary.Blocks)-1 {
		//		curPos += fadeWidth * 0.5
		//	}
		//	if fadeDuration > 0 {
		//		pushFrame()
		//	}
		//	lastTimeStamp += fadeDuration
		//}
		//fmt.Println(frames)
		//anim := word.Ele.Call("animate", frames, map[string]interface{}{
		//	"duration": 1000,
		//	"easing":   "linear",
		//	"fill":     "both",
		//})
		//anim.Call("play")
		//fmt.Println(totalFadeDuration)
		//if anim.Truthy() {
		//	progress := currentTime - word.Begin
		//	if progress >= 0 && progress <= (word.End-word.Begin) {
		//		fmt.Println(progress)
		//		anim.Set("currentTime", float64(progress))
		//		anim.Call("play")
		//	} else {
		//		anim.Call("pause")
		//	}
		//}

		if word.Begin == 0 && word.End == 0 {
			continue
		}
		intervalTime := word.Begin - currentTime

		// 计算动画参数
		duration := word.End - word.Begin
		//eleWidth := word.Ele.Get("offsetWidth").Float()
		//oPeo := float64(bglw) + eleWidth
		//oldV := eleWidth / duration.Seconds()
		//newT := (eleWidth + float64(bglw)) / oldV
		//newV := (eleWidth + float64(bglw)) / newT
		//nT := (eleWidth / newV) + float64(bglw)/newV
		// 设置 GSAP 动画
		//animation := gsap.Call("to", word.Ele, map[string]interface{}{
		//	"duration": duration.Seconds(),
		//	//"duration": 0,
		//	"ease": "none",
		//	//"--p":   fmt.Sprintf("%vpx", oPeo),
		//	"background-position-x": "0px",
		//	"delay":                 (intervalTime.Seconds()),
		//})
		//
		//if wi+1 < len(words) {
		//	frx := getFPX(word.Ele.Get("offsetWidth").Float(), word.Ele.Get("offsetHeight").Float(), fadeRatio)
		//	v := frx / duration.Seconds()
		//	t := words[wi+1].Ele.Get("offsetWidth").Float() / v
		//	fmt.Println(t)
		//	//intTime := words[wi+1].Begin - currentTime
		//	gsap.Call("to", words[wi+1].Ele, map[string]interface{}{
		//		"duration": t,
		//		//"duration": 0,
		//		"ease": "none",
		//		//"--p":   fmt.Sprintf("%vpx", oPeo),
		//		"background-position-x": fmt.Sprintf("%vpx", -words[wi+1].Ele.Get("offsetWidth").Float()),
		//		"delay":                 (intervalTime.Seconds() + duration.Seconds() - t),
		//	})
		//}

		// 计算时间
		//startOffset := word.Begin - lrcs.Contents[index].Primary.Begin

		//wordWidth := word.Ele.Get("offsetWidth").Float()
		//height := word.Ele.Get("offsetHeight").Float()
		// 渐变宽度（推荐统一计算）
		//width := wordWidth // 已测量好的宽度
		///fadeW := height * fadeRatio
		//totalMove := width + fadeW*2

		// 计算动画起止位置
		//fromX := -curX - totalMove
		//toX := -curX

		// 动画
		//animation := gsap.Call("to", word.Ele,
		//	map[string]interface{}{
		//		"background-position-x": "0px",
		//		"duration":              duration.Seconds(),
		//		"delay":                 startOffset.Seconds(),
		//		"ease":                  "none",
		//	})
		//
		//curX += width
		//animation.Call("eventCallback", "onComplete", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//	// word.AheadtrAnimation是否为空
		//	if !word.AheadtrAnimation.IsUndefined() && !word.AheadtrAnimation.IsNull() {
		//		fmt.Println("kill Animation")
		//		word.AheadtrAnimation.Call("pause")
		//		word.AheadtrAnimation.Call("kill")
		//		// 让animation的进度完成
		//		animation.Call("kill")
		//	}
		//	return nil
		//}))
		//fmt.Println((intervalTime.Seconds()))
		//if wi+1 < len(words) {
		//	word.AheadtrAnimation = gsap.Call("to", words[wi+1].Ele, map[string]interface{}{
		//		"duration": getLastOffsetTime(wi-1, lrcs.Contents[index].Primary).Seconds(),
		//		"ease":     "none",
		//		"--p":      fmt.Sprintf("%vpx", bglw),
		//		"delay":    (intervalTime.Seconds() + (eleWidth / newV)),
		//	})
		//
		//}
		//word.GsapAnimation = animation

		//GsetTimeout(func() {

		//fst := time.Duration(0)
		//f wi+1 < len(words) {
		//	fmt.Println(word.Text, "----------------", words[wi+1].Text)
		//	d := getLastOffsetTime(wi+1, lrcs.Contents[index].Primary)
		//	ofst += d
		//	fmt.Println(d.Seconds())
		//	intTime := words[wi].End - currentTime
		//	words[wi+1].Ele.Call("animate", []interface{}{
		//		map[string]interface{}{
		//			"backgroundPositionX": fmt.Sprintf("%vpx", -words[wi+1].Ele.Get("offsetWidth").Float()),
		//		}},
		//		map[string]interface{}{
		//			"duration": d.Milliseconds(),
		//			"easing":   "linear",
		//			"delay":    float64(intTime.Milliseconds()) - (word.Ele.Get("offsetHeight").Float()*fadeRatio)/(word.Ele.Get("offsetHeight").Float()/float64(duration.Milliseconds())),
		//			"fill":     "forwards",
		//		},
		//	)
		//
		//nimation := word.Ele.Call("animate", []interface{}{
		//	map[string]interface{}{
		//		"backgroundPositionX": "0px",
		//	}},
		//	map[string]interface{}{
		//		"duration": duration.Milliseconds(),
		//		"easing":   "linear",
		//		"delay":    float64(intervalTime.Milliseconds()),
		//		"fill":     "forwards",
		//	},
		//

		//ord.GsapAnimation = animation

		frame := createFrames(wordsN, wi, time.Duration(lineTime))
		log.Println("linetime", lineTime)
		word.GsapAnimation = word.Ele.Call("animate", js.ValueOf(frame), map[string]interface{}{
			"duration": time.Duration(lineTime).Milliseconds(),
			"easing":   "linear",
			"fill":     "forwards",
		})

		upAnimateTime := duration.Milliseconds() + 700
		aimt := word.Ele.Call("animate", []interface{}{
			map[string]interface{}{},
			map[string]interface{}{
				"transform": "translateY(5px)",
				//"marginTop": "-5px",
			},
		},
			map[string]interface{}{
				"duration": upAnimateTime,
				"easing":   "ease-out",
				//"delay":    (intervalTime.Milliseconds() * 95 / 100),
				"delay": float64(intervalTime.Milliseconds()),
				"fill":  "forwards",
			},
		)

		var bgChildren js.Value = word.Ele.Get("children")
		letterDuration := float64(duration.Milliseconds()) / (float64(len(word.Text)) - (float64(len(word.Text))-1)*0.7)
		chrDu := duration.Seconds() / float64(len(word.Text))

		for i := 0; i < bgChildren.Length(); i++ {

			item := bgChildren.Index(i)

			var charWidth = item.Get("offsetWidth").Float()
			var oldV = charWidth / chrDu
			var ope = float64(bglw) + charWidth
			var nT = ope / oldV
			hlga := gsap.Call("to", item, map[string]interface{}{
				// duration单位为秒
				"duration": nT,
				"ease":     "none",
				"--p":      fmt.Sprintf("%vpx", ope),
				// 延时触发
				"delay": float64(i)*chrDu + intervalTime.Seconds(),
			})
			word.HighLightBackgroungAnimation = append(word.HighLightBackgroungAnimation, hlga)

			//const letterDuration = totalDuration / (letterCount - (letterCount - 1) * overlapRatio);
			//const startTime = index * letterDuration * (1 - overlapRatio);
			startTime := float64(i) * letterDuration * (1 - 0.7)
			aimt := item.Call("animate", []interface{}{
				map[string]interface{}{
					"easing": "ease",
				},
				map[string]interface{}{
					//"transform": "translateY(-3.6%)",
					//transform: "matrix3d(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, -10, 0, 0.8)",
					//"textShadow": "rgba(255, 255, 255, 0.20) 0px 0px " + "20" + "px",
					"transform": "scale(1.15) translateX(" + fmt.Sprintf("%f", getScaleOffset(i, 1.15, word.Ele)) + "px) translateY(1%)",
					"easing":    "ease",
					//"color":     "rgba(255, 255, 255, 1)",
					"filter": "blur(0.5px)",
				},
				//map[string]interface{}{
				//	"transform":  "translateX(0px) scale(1)",
				//	"textShadow": "none",
				//	"easing":     "ease",
				//	"color":      "rgba(255, 255, 255,1)",
				//	"filter":     "blur(0px)",
				//},
			},

				map[string]interface{}{
					"duration": letterDuration * 1.5,
					"fill":     "forwards",
					//"delay":    float64((i)*4)/10*float64(duration.Milliseconds())/float64(bgChildren.Length()) + float64(intervalTime.Milliseconds()*95/100),
					//"delay": float64(i)*float64(duration.Milliseconds())*0.2 - float64(duration.Milliseconds())*0.1*float64(i) + float64(intervalTime.Milliseconds()*95/100),
					"delay": startTime + float64(intervalTime.Milliseconds()*95/100),
				},
			)

			item.Call("animate", []interface{}{
				map[string]interface{}{
					"easing": "ease",
				},
				map[string]interface{}{
					"transform":  "translateX(0px) scale(1) translateY(0)",
					"textShadow": "none",
					"easing":     "ease",
					//"color":      "rgba(255, 255, 255,1)",
					"filter": "blur(0px)",
				},
			},

				map[string]interface{}{
					"duration": duration.Milliseconds(),
					"fill":     "forwards",
					//"delay":    float64((i)*4)/10*float64(duration.Milliseconds())/float64(bgChildren.Length()) + float64(intervalTime.Milliseconds()*95/100),
					"delay": float64(duration.Milliseconds()) + float64(intervalTime.Milliseconds()*95/100),
				},
			)
			//if i == bgChildren.Length()-1 {
			//	aimt.Call("addEventListener", "finish", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			//		for j := 0; j < bgChildren.Length(); j++ {
			//			item := bgChildren.Index(j)
			//			item.Call("animate", []interface{}{
			//				map[string]interface{}{},
			//				map[string]interface{}{
			//					"transform": "scale(1) translateY(0px)",
			//					//transform: "matrix(1, 0, 0, 1, 0, 0)",
			//					"textShadow": "none",
			//					"easing":     "cubic-bezier(0.5, 0, 0.5, 1)",
			//					"color":      "rgba(255, 255, 255,1)",
			//					"filter":     "blur(0px)",
			//					//marginLeft: "none"
			//				},
			//			},
			//
			//				map[string]interface{}{
			//					"duration": 700,
			//					"easing":   "ease",
			//					"fill":     "forwards",
			//				},
			//			)
			//		}
			//		return nil
			//	}))
			//}
			word.HighLightAnimations = append(word.HighLightAnimations, aimt)
		}

		//aimt.Call("addEventListener", "finish", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//	word.Ele.Get("style").Set("transform", "translateY(5px)")
		//	return nil
		//}))
		word.TextUpAnimation = aimt
	}

}

/*
function getScaleOffset(index, scale) {
            const centerIndex = (chars.length - 1) / 2;
            const baseWidth = chars[0].offsetWidth;
            return (index - centerIndex) * baseWidth * (scale - 1) * 0.5;
        }*/

func getScaleOffset(index int, scale float64, dom js.Value) float64 {
	chars := dom.Get("children")
	centerIndex := (chars.Length() - 1) / 2
	baseWidth := chars.Index(0).Get("offsetWidth").Float()
	return float64(index-centerIndex) * baseWidth * (scale - 1) * 0.5
}

func getLastOffsetTime(index int, line *lyrics.Line) time.Duration {
	// 默认返回值
	defaultDuration := 500 * time.Millisecond

	// 反向遍历
	for i := index - 1; i >= 0; i-- {
		block := line.Blocks[i]
		if block.Begin != 0 && block.End != 0 {
			fmt.Println("找到了:", line.Blocks[index].Text, "上一个为:", line.Blocks[i].Text)
			offsetWidth := block.Ele.Get("offsetWidth").Float()
			//
			//// 持续时间（直接使用毫秒）
			fmt.Println("end:", block.End, " - begin:", block.Begin)
			durationMs := float64(block.End.Milliseconds() - block.Begin.Milliseconds())
			fmt.Println(block.Text, " 的持续时间为：（ms）", durationMs)
			//
			//// 避免除以0
			//if durationMs == 0 {
			//	return defaultDuration
			//}
			//
			//// 计算速度（像素/毫秒）
			speed := offsetWidth / durationMs
			fmt.Println(block.Text, " 的速度为：（px/ms）", speed)
			//
			//// 避免速度为0
			if speed == 0 {
				return defaultDuration
			}

			//// 计算额外时间（毫秒）
			extraTimeMs := ((line.Blocks[index].Ele.Get("offsetHeight").Float() * fadeRatio) / speed) * float64(time.Millisecond)
			fmt.Println(line.Blocks[index].Text, " 运行", (line.Blocks[i].Ele.Get("offsetHeight").Float() * fadeRatio), "ps的时间为：（ms）", extraTimeMs)
			//
			//// 将毫秒转为Duration（更安全的转换方式）
			return time.Duration(extraTimeMs)
			//return defaultDuration
		}
	}

	return defaultDuration
}
func generateBackgroundFadeStyle(elementWidth, elementHeight, fadeRatio float64) (string, string, float64, float64) {

	/*const fadeWidth = elementHeight * fadeRatio;
	  const widthRatio = fadeWidth / elementWidth;

	  // 使用源码同款算法
	  const totalAspect = 2 + widthRatio;
	  const widthInTotal = widthRatio / totalAspect;
	  const leftPos = (1 - widthInTotal) / 2;

	  const from = leftPos * 100;
	  const to = (leftPos + widthInTotal) * 100;

	  const backgroundImage = `linear-gradient(to right,
	    rgba(255, 255, 255, 1.0) ${from.toFixed(6)}%,
	    rgba(255, 255, 255, 0.0) ${to.toFixed(6)}%)`;

	  const backgroundSize = `${(totalAspect * 100).toFixed(3)}% 100%`;
	  const totalPxWidth = elementWidth + fadeWidth;

	  return {
	    backgroundImage,
	    backgroundSize,
	    backgroundRepeat: 'no-repeat',
	    backgroundOrigin: 'left',
	    backgroundPositionX: `${-totalPxWidth}px`,
	    finalPositionX: `0px`,
	    transitionDistance: totalPxWidth,
	  };*/

	fadeWidth := elementHeight * fadeRatio
	widthRatio := fadeWidth / elementWidth

	totalAspect := 2 + widthRatio
	widthInTotal := widthRatio / totalAspect
	leftPos := (1 - widthInTotal) / 2

	from := leftPos * 100
	to := (leftPos + widthInTotal) * 100

	backgroundImage := fmt.Sprintf("linear-gradient(to right,rgba(255, 255, 255, var(--rcolor)) %f%%,rgba(255, 255, 255, var(--color)) %f%%)", from, to)
	backgroundSize := fmt.Sprintf("%f%% 100%%", totalAspect*100)
	totalPxWidth := elementWidth + fadeWidth
	return backgroundImage, backgroundSize, -totalPxWidth, totalPxWidth
}

func getFPX(width, height, fr float64) float64 {
	return height*fr + width
}

/*func createFraems(blocks []*lyrics.Block, index int, lineTime time.Duration) []interface{} {
	var frames []interface{}
	ElWidth := blocks[index].Ele.Get("offsetWidth").Float()
	ElHeight := blocks[index].Ele.Get("offsetHeight").Float()
	hr := ElHeight * fadeRatio
	fbw := ElWidth + hr

	lastWordIndex := index - 1
	if lastWordIndex < 0 {
		lastWordIndex = 0

		frames = append(frames, map[string]interface{}{
			"backgroundPositionX": fmt.Sprintf("-%fpx", fbw),
			//"--translateX": fmt.Sprintf("%fpx", 0),
			"offset": 0,
		})

		frames = append(frames, map[string]interface{}{
			"backgroundPositionX": fmt.Sprintf("-%fpx", fbw),

			//"--translateX": fmt.Sprintf("%fpx", 0),
			"offset": float64(blocks[index].Begin.Milliseconds()-blocks[0].Begin.Milliseconds()) / float64(lineTime.Milliseconds()),
		})
		frames = append(frames, map[string]interface{}{
			"backgroundPositionX": "0px",
			//"--translateX": fmt.Sprintf("%fpx", ElWidth),
			"offset": float64(blocks[index].End.Milliseconds()-blocks[0].Begin.Milliseconds()) / float64(lineTime.Milliseconds()),
		})

		frames = append(frames, map[string]interface{}{
			"backgroundPositionX": "0px",
			//"--translateX": fmt.Sprintf("%fpx", ElWidth),
			"offset": 1,
		})
		fmt.Println(frames)
		return frames
	}

	lastWordElWidth := blocks[lastWordIndex].Ele.Get("offsetWidth").Float()
	lastWordElHeight := blocks[lastWordIndex].Ele.Get("offsetHeight").Float()
	lastWordHr := lastWordElHeight * fadeRatio
	//lastWordFBW := lastWordElWidth + lastWordHr
	lastWordDuration := blocks[lastWordIndex].End.Milliseconds() - blocks[lastWordIndex].Begin.Milliseconds()
	lastWordV := lastWordElWidth / float64(lastWordDuration)
	lastWordHrTime := lastWordHr / lastWordV
	//lastWordExTime:= lastWordDuration - int64(lastWordHrTime)
	fmt.Println("end:", blocks[index].End.Milliseconds(), " - begin:", blocks[index].Begin.Milliseconds(), " - lastWordHrTime:", lastWordHrTime)
	frames = append(frames, map[string]interface{}{
		"backgroundPositionX": fmt.Sprintf("-%fpx", fbw),
		//"--translateX": fmt.Sprintf("%fpx", 0),
		"offset": 0,
	})

	frames = append(frames, map[string]interface{}{
		"backgroundPositionX": fmt.Sprintf("-%fpx", fbw),

		//"--translateX": fmt.Sprintf("%fpx", 0),
		"offset": float64(blocks[lastWordIndex].End.Milliseconds()-blocks[0].Begin.Milliseconds()-int64(lastWordHrTime)) / float64(lineTime.Milliseconds()),
	})

	frames = append(frames, map[string]interface{}{
		"backgroundPositionX": fmt.Sprintf("-%fpx", fbw-hr),

		//"--translateX": fmt.Sprintf("%fpx", 0),
		"offset": float64(blocks[lastWordIndex].End.Milliseconds()-blocks[0].Begin.Milliseconds()) / float64(lineTime.Milliseconds()),
	})

	frames = append(frames, map[string]interface{}{
		"backgroundPositionX": fmt.Sprintf("-%fpx", fbw-hr),

		//"--translateX": fmt.Sprintf("%fpx", 0),
		"offset": float64(blocks[index].Begin.Milliseconds()-blocks[0].Begin.Milliseconds()) / float64(lineTime.Milliseconds()),
	})
	frames = append(frames, map[string]interface{}{
		"backgroundPositionX": "0px",
		//"--translateX": fmt.Sprintf("%fpx", ElWidth),
		"offset": float64(blocks[index].End.Milliseconds()-blocks[0].Begin.Milliseconds()) / float64(lineTime.Milliseconds()),
	})

	frames = append(frames, map[string]interface{}{
		"backgroundPositionX": "0px",
		//"--translateX": fmt.Sprintf("%fpx", ElWidth),
		"offset": 1,
	})
	fmt.Println(frames)
	return frames
}
*/

func createFrames(blocks []*lyrics.Block, index int, lineTime time.Duration) []interface{} {
	var frames []interface{}
	ElWidth := blocks[index].Ele.Get("offsetWidth").Float()
	ElHeight := blocks[index].Ele.Get("offsetHeight").Float()
	hr := ElHeight * fadeRatio
	fbw := ElWidth + hr

	// 计算总持续时间(以最后一个单词的结束时间为准)
	totalDuration := lineTime
	if len(blocks) > 0 {
		lastBlockEnd := blocks[len(blocks)-1].End
		totalDuration = lastBlockEnd - blocks[0].Begin
	}

	// 计算当前单词之前的累计宽度
	widthBeforeSelf := 0.0
	for i := 0; i < index; i++ {
		widthBeforeSelf += blocks[i].Ele.Get("offsetWidth").Float()
	}
	if index > 0 {
		widthBeforeSelf += hr // 第一个单词有额外的渐变宽度
	}

	minOffset := -fbw
	clampOffset := func(x float64) float64 {
		if x < minOffset {
			return minOffset
		}
		if x > 0 {
			return 0
		}
		return x
	}

	currentPos := -widthBeforeSelf - ElWidth - hr
	timeOffset := 0.0
	lastPos := currentPos
	lastTime := 0.0

	pushFrame := func() {
		// 确保时间在0-1范围内
		time := math.Max(0, math.Min(1, timeOffset))
		duration := time - lastTime
		moveOffset := currentPos - lastPos
		d := 0.0
		if moveOffset != 0 {
			d = math.Abs(duration / moveOffset)
		}

		// 处理边界情况
		if currentPos > minOffset && lastPos < minOffset {
			staticTime := math.Abs(lastPos-minOffset) * d
			frames = append(frames, map[string]interface{}{
				"backgroundPositionX": fmt.Sprintf("%fpx", clampOffset(lastPos)),
				"offset":              lastTime + staticTime,
			})
		}

		if currentPos > 0 && lastPos < 0 {
			staticTime := math.Abs(lastPos) * d
			frames = append(frames, map[string]interface{}{
				"backgroundPositionX": fmt.Sprintf("%fpx", clampOffset(currentPos)),
				"offset":              lastTime + staticTime,
			})
		}

		frames = append(frames, map[string]interface{}{
			"backgroundPositionX": fmt.Sprintf("%fpx", clampOffset(currentPos)),
			"offset":              time,
		})

		lastPos = currentPos
		lastTime = time
	}

	// 初始帧
	pushFrame()

	lastTimeStamp := 0.0
	for i, block := range blocks {
		// 停顿阶段
		curTimeStamp := float64((block.Begin - blocks[0].Begin).Milliseconds())
		staticDuration := curTimeStamp - lastTimeStamp
		timeOffset += staticDuration / float64(totalDuration.Milliseconds())
		if staticDuration > 0 {
			pushFrame()
		}
		lastTimeStamp = curTimeStamp

		// 移动阶段
		fadeDuration := float64((block.End - block.Begin).Milliseconds())
		timeOffset += fadeDuration / float64(totalDuration.Milliseconds())
		currentPos += block.Ele.Get("offsetWidth").Float()

		// 第一个和最后一个单词有额外的渐变处理
		if i == 0 {
			currentPos += hr * 1.5
		}
		if i == len(blocks)-1 {
			currentPos += hr * 0.5
		}

		if fadeDuration > 0 {
			pushFrame()
		}
		lastTimeStamp += fadeDuration
	}

	return frames
}
