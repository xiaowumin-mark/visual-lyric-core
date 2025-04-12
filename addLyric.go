package main

import (
	"fmt"
	"syscall/js"
	"time"

	lyrics "github.com/xiaowumin-mark/visual-lyric-core/lyric"
)

func addLyric(index int, lrcs *lyrics.Lyrics) {

	addIndex(index)
	// 在js控制台打印nowPlayingIndex
	currentTime := getCurrentTime(audio)
	words := lrcs.Contents[index].Primary.Blocks

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
				animation := gsap.Call("to", word.Ele, map[string]interface{}{
					// duration单位为秒
					"duration": duration.Seconds() * 1.05,
					"ease":     "none",
					"--p":      "100%",
					"--rp":     "140%",
					// 延时触发
					"delay": intervalTime.Seconds() * 0.95,
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
	var wLSTs []float64
	for _, word := range words {
		if word.Begin == 0 && word.End == 0 {
			continue
		}
		intervalTime := word.Begin - currentTime

		// 计算动画参数
		duration := word.End - word.Begin
		eleWidth := word.Ele.Get("offsetWidth").Float()

		// 调整后的动画终点（从 -20px 到 100%）
		oPeo := (float64(bglw) * 2 / eleWidth * 100) + 100

		V := 100 / duration.Seconds() // 原速度（100% / duration）

		// 计算 20px 过渡的额外时间
		extraTime := (float64(bglw) * 2 / eleWidth * 100) / (oPeo / duration.Seconds())
		wLSTs = append(wLSTs, extraTime)

		// 计算累计提前时间
		var ofst float64 = 0
		for ind, ti := range wLSTs {
			if ind == 0 {
				ofst += 0

			} else {
				ofst += ti
			}

		}

		// 调整后的动画时间（保持速度 V 不变）
		nT := oPeo / V

		// 设置 GSAP 动画
		animation := gsap.Call("to", word.Ele, map[string]interface{}{
			"duration": nT * 1.1,
			//"duration": 0,
			"ease":  "none",
			"--p":   fmt.Sprintf("%v%%", oPeo),
			"delay": (intervalTime.Seconds() - extraTime) * 0.95,
		})
		word.GsapAnimation = animation
		//GsetTimeout(func() {
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
				"delay": (float64(intervalTime.Milliseconds())-float64(extraTime*1000))*0.95 + 200,
				"fill":  "forwards",
			},
		)

		var bgChildren js.Value = word.Ele.Get("children")
		letterDuration := float64(duration.Milliseconds()) / (float64(len(word.Text)) - (float64(len(word.Text))-1)*0.7)
		chrDu := duration.Seconds() / float64(len(word.Text))
		for i := 0; i < bgChildren.Length(); i++ {

			item := bgChildren.Index(i)

			var charWidth = item.Get("offsetWidth").Float()
			var oldV = 100 / (duration.Seconds() / float64(len(word.Text)))
			var ope = (float64(bglw) * 2 / charWidth * 100) + 100
			var nT = ope / oldV
			hlga := gsap.Call("to", item, map[string]interface{}{
				// duration单位为秒
				"duration": nT,
				"ease":     "none",
				"--p":      fmt.Sprintf("%v%%", ope),
				// 延时触发
				"delay": float64(i)*chrDu + float64(intervalTime.Seconds()*95/100),
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
					//"textShadow": "rgba(255, 255, 255, 0.70) 0px 0px " + "15" + "px",
					"transform": "scale(1.1) translateX(" + fmt.Sprintf("%f", getScaleOffset(i, 1.1, word.Ele)) + "px) translateY(1%)",
					"easing":    "ease",
					//"color":     "rgba(255, 255, 255, 1)",
					"filter": "blur(1.25px)",
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
