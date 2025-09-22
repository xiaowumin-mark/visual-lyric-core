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
	lineTime := lrcs.Contents[index].Primary.Blocks[len(lrcs.Contents[index].Primary.Blocks)-1].End - lrcs.Contents[index].Primary.Begin
	if lineTime < 0 {
		lineTime = lrcs.Contents[index].Primary.End - lrcs.Contents[index].Primary.Begin
	}
	fmt.Println(time.Duration(lineTime).Milliseconds())

	gsap.Call("to", lrcs.Contents[index].Ele, map[string]interface{}{
		"duration": 0.5,
		"scale":    1.05,
		"delay":    0.1,
	})

	if len(lrcs.Contents[index].Backgrounds) != 0 {
		//lrcs.Contents[index].BackgroundsEle.Get("style").Set("display", "block")
		//GsetTimeout(func() {
		//lrcs.Contents[index].BackgroundsEle.Get("classList").Call("add", "bgShow")
		lrcs.Contents[index].ShowBackgrounds = true
		GsetTimeout(func() {
			gd(bubbleSort(nowPlayingIndex)[0], lrcs, false)
		}, 50*time.Millisecond)
		//}, 10*time.Millisecond)

		for _, item := range lrcs.Contents[index].Backgrounds {

			//item.Ele.Call("animate", []interface{}{
			//	map[string]interface{}{
			//		"opacity":   0,
			//		"transform": "scale(0.8) " + item.Ele.Get("style").Get("transform").String(),
			//	},
			//	map[string]interface{}{
			//		"opacity":   1,
			//		"transform": "scale(1) " + item.Ele.Get("style").Get("transform").String(),
			//	},
			//}, map[string]interface{}{
			//	"duration": 500,
			//	"fill":     "forwards",
			//	"easing":   "ease-out",
			//})
			gsap.Call("to", item.Ele, map[string]interface{}{
				"opacity":  1,
				"scale":    1,
				"duration": 0.5,
				"delay":    0.2,
			})

			bgLineTime := item.Blocks[len(item.Blocks)-1].End - item.Blocks[0].Begin
			delay := item.Blocks[0].Begin - lrcs.Contents[index].Primary.Begin
			var bgWordsN []*lyrics.Block
			for _, word := range item.Blocks {
				backgroundImage, backgroundSize, backgroungPX, _ := generateBackgroundFadeStyle(word.Ele.Get("offsetWidth").Float(), word.Ele.Get("offsetHeight").Float(), bgfadeRatio)
				word.Ele.Get("style").Set("backgroundImage", backgroundImage)
				word.Ele.Get("style").Set("backgroundSize", backgroundSize)
				word.Ele.Get("style").Set("backgroundPositionX", fmt.Sprintf("%vpx", backgroungPX))
				if word.Begin == 0 && word.End == 0 {
					continue
				}
				bgWordsN = append(bgWordsN, word)
			}
			for wi, word := range bgWordsN {
				frame := createFrames(bgWordsN, wi, bgLineTime, bgfadeRatio)
				animate := word.Ele.Call("animate", js.ValueOf(frame), map[string]interface{}{
					"duration": bgLineTime.Milliseconds(),
					"easing":   "linear",
					"fill":     "forwards",
					"delay":    delay.Milliseconds(),
				})
				word.Animation = append(word.Animation, animate)

				intervalTime := word.Begin - currentTime

				// 计算动画参数
				duration := word.End - word.Begin
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
				word.TextUpAnimation = aimt
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
	var wordsWS []*lyrics.Block

	lrcs.Contents[index].Ele.Get("style").Set("filter", "blur(0px)")

	for _, word := range words {
		if word.Begin == 0 && word.End == 0 {
			continue
		}
		wordsN = append(wordsN, word)

		var bgChildren js.Value = word.Ele.Get("children")
		if bgChildren.Length() > 0 {
			wsitemTime := (word.End - word.Begin) / time.Duration(bgChildren.Length())
			for i := 0; i < bgChildren.Length(); i++ {
				var bgChild js.Value = bgChildren.Index(i)
				wordsWS = append(wordsWS, &lyrics.Block{
					Ele:   bgChild,
					Text:  bgChild.Get("innerHTML").String(),
					Begin: time.Duration(word.Begin) + time.Duration(i)*wsitemTime,
					End:   time.Duration(word.Begin) + time.Duration(i+1)*wsitemTime,
				})
			}
		} else {
			wordsWS = append(wordsWS, word)
		}
	}
	var lineWordsAnimates []js.Value
	//curX := 0.0 // 累积偏移量（整行动画进度）
	log.Println("lineTime:", lineTime)
	for wi, word := range wordsWS {
		log.Println(wordsN[len(wordsN)-1].End.Seconds(), wordsN[0].Begin.Seconds())
		frame := createFrames(wordsWS, wi, wordsN[len(wordsN)-1].End-wordsN[0].Begin, fadeRatio)
		wordAimate := word.Ele.Call("animate", js.ValueOf(frame), map[string]interface{}{
			"duration": time.Duration(lineTime).Milliseconds(),
			"easing":   "linear",
			"fill":     "forwards",
		})
		//word.Animation = append(word.Animation, wordAimate)
		lineWordsAnimates = append(lineWordsAnimates, wordAimate)
	}

	lastWordWSIndex := 0
	// 将lineWordsAnimates和wordsN绑定
	for _, wordItem := range wordsN {
		var bgChildren js.Value = wordItem.Ele.Get("children")
		if bgChildren.Length() > 0 {
			wordItem.Animation = lineWordsAnimates[lastWordWSIndex : lastWordWSIndex+bgChildren.Length()]
			lastWordWSIndex += bgChildren.Length()
		} else {
			wordItem.Animation = append(wordItem.Animation, lineWordsAnimates[lastWordWSIndex])
			lastWordWSIndex += 1
		}

	}
	//for _, word := range wordsN {

	//if word.Begin == 0 && word.End == 0 {
	//	continue
	//}

	/*
		var bgChildren js.Value = word.Ele.Get("children")

		for i := 0; i < bgChildren.Length(); i++ {
			pjsj := duration / time.Duration(bgChildren.Length())
			baseDelay := word.Begin - currentTime
			charDelay := pjsj * time.Duration(i)

			aimt := bgChildren.Index(i).Call("animate", []interface{}{
				map[string]interface{}{

					"transform": "scale(1)",
					"easing":    "ease-out",
				},
				map[string]interface{}{
					"transform": "scale(1.15) translateX(" + fmt.Sprintf("%f", getScaleOffset(i, 1.15, word.Ele)) + "px) translateY(1%)",
					"easing":    "ease-in",
				},
				map[string]interface{}{
					"transform": "scale(1)",
					"easing":    "ease",
				},
			},
				map[string]interface{}{
					"duration": (duration * 3 / 2).Milliseconds(),
					"delay":    (baseDelay + charDelay - charDelay*60/100).Milliseconds(), // 通过系数控制重叠比例
					"fill":     "forwards",
					"easing":   "ease",
				},
			)
			word.HighLightAnimations = append(word.HighLightAnimations, aimt)
		}*/

	//}

	for _, word := range lrcs.Contents[index].Primary.WrdList {

		item := lrcs.Contents[index]
		//startTime := item.Primary.Blocks[word[0]].Begin
		//endTime := item.Primary.Blocks[word[len(word)-1]].End
		var duration time.Duration
		var allWordEles []js.Value
		for i := word[0]; i <= word[len(word)-1]; i++ {
			var blanks = item.Primary.Blocks[i].Ele.Get("children")
			for j := 0; j < blanks.Length(); j++ {
				allWordEles = append(allWordEles, blanks.Index(j))
			}

			wordDuration := item.Primary.Blocks[i].End - item.Primary.Blocks[i].Begin
			duration += wordDuration
		}
		litIndex := 0
		wordIndex := 0
		for i := word[0]; i <= word[len(word)-1]; i++ {
			var blanks = item.Primary.Blocks[i].Ele.Get("children")
			//intervalTime := word.Begin - currentTime

			// 计算动画参数
			upAnimateTime := duration.Milliseconds() + 700
			aimt := item.Primary.Blocks[i].Ele.Call("animate", []interface{}{
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
					"delay": float64((item.Primary.Blocks[word[0]].Begin - currentTime).Milliseconds()),
					"fill":  "forwards",
				},
			)
			item.Primary.Blocks[i].TextUpAnimation = aimt

			wordDuration := item.Primary.Blocks[i].End - item.Primary.Blocks[i].Begin

			for j := 0; j < blanks.Length(); j++ {
				pjsj := wordDuration / time.Duration(blanks.Length())
				baseDelay := item.Primary.Blocks[i].Begin - currentTime
				charDelay := pjsj * time.Duration(j)

				aimt := blanks.Index(j).Call("animate", []interface{}{
					map[string]interface{}{

						"transform": "scale(1)",
						"easing":    "ease-out",
					},
					map[string]interface{}{
						"transform": "scale(1.1) translateX(" + fmt.Sprintf("%f", getScaleOffset(litIndex, 1.1, allWordEles)) + "px) translateY(1%)",
						"easing":    "ease-in",
					},
					map[string]interface{}{
						"transform": "scale(1)",
						"easing":    "ease",
					},
				},
					map[string]interface{}{
						"duration": (duration * 3 / 2).Milliseconds(),
						"delay":    (baseDelay + charDelay - charDelay*60/100).Milliseconds(), // 通过系数控制重叠比例
						"fill":     "forwards",
						"easing":   "ease",
					},
				)
				item.Primary.Blocks[i].HighLightAnimations = append(item.Primary.Blocks[i].HighLightAnimations, aimt)
				litIndex++
			}
			wordIndex++
		}

	}

}

func getScaleOffset(index int, scale float64, doms []js.Value) float64 {
	//chars := dom.Get("children")
	centerIndex := (len(doms) - 1) / 2

	/*centerIndex := (chars.Length() - 1) / 2
	baseWidth := chars.Index(0).Get("offsetWidth").Float()
	return float64(index-centerIndex) * baseWidth * (scale - 1) * 0.5*/

	// Calculate the cumulative width up to the current character
	cumulativeWidth := 0.0
	for i := 0; i < index; i++ {
		//cumulativeWidth += doms[i].offsetWidth;
		cumulativeWidth += doms[i].Get("offsetWidth").Float()
	}

	// Calculate the cumulative width up to the center character
	centerCumulativeWidth := 0.0
	for i := 0; i < centerIndex; i++ {
		//centerCumulativeWidth += doms[i].offsetWidth;
		centerCumulativeWidth += doms[i].Get("offsetWidth").Float()
	}

	// The offset is the difference between current position and center position,
	// multiplied by the scale factor
	return (cumulativeWidth - centerCumulativeWidth) * (scale - 1) * 0.5
}

func generateBackgroundFadeStyle(elementWidth, elementHeight, fadeRatio float64) (string, string, float64, float64) {

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

func createFrames(blocks []*lyrics.Block, index int, lineTime time.Duration, fadeRatio float64) []interface{} {
	var frames []interface{}

	/*backgroundImage, backgroundSize, backgroungPX, _ := generateBackgroundFadeStyle(blocks[index].Ele.Get("offsetWidth").Float(), blocks[index].Ele.Get("offsetHeight").Float(), fadeRatio)
	blocks[index].Ele.Get("style").Set("backgroundImage", backgroundImage)
	blocks[index].Ele.Get("style").Set("backgroundSize", backgroundSize)
	blocks[index].Ele.Get("style").Set("backgroundPositionX", fmt.Sprintf("%vpx", backgroungPX))*/

	ElWidth := blocks[index].Ele.Get("offsetWidth").Float()
	ElHeight := blocks[index].Ele.Get("offsetHeight").Float()
	hr := ElHeight * fadeRatio
	fbw := ElWidth + hr - 2

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
	if len(frames) > 0 {
		if lastFrame, ok := frames[len(frames)-1].(map[string]interface{}); ok {
			lastFrame["offset"] = 1.0
		}
	}
	log.Println(frames...)
	return frames
}
